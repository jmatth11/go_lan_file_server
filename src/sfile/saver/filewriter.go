package saver

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"math"
	"os"
	"sfile"
	"sfile/conversion"
)

const (
	headerAppend = "-Header"
	blockSize    = 2000000      // 2MB
	maxFileSize  = 128000000000 // 128GB but can change if expecting large video files
)

// FileWriter object for creating/updating/deleting/reading a SaveFile object on a file system.
type FileWriter struct {
	dataFile   *os.File
	headerFile *os.File
	path       string
}

// NewFileWriter takes a path to the SaveFile and returns a new pointed FileWriter
func NewFileWriter(path string) (*FileWriter, error) {
	return &FileWriter{
		path: path,
	}, nil
}

// grabFilesForWrite creates the main data file and the associated header file.
func (fw *FileWriter) grabFilesForWrite(fileName []byte) error {
	nameBuffer := bytes.NewBufferString(fw.path)
	nameBuffer.WriteRune(os.PathSeparator)
	nameBuffer.Write(fileName)
	file1, err := os.OpenFile(nameBuffer.String(), os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	nameBuffer.WriteString(headerAppend)
	file2, err := os.OpenFile(nameBuffer.String(), os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	fw.dataFile = file1
	fw.headerFile = file2
	return nil
}

// Create takes a SaveFile object and creates the header and data files needed
func (fw *FileWriter) Create(data sfile.DataFormat) error {
	if err := fw.grabFilesForWrite(data.FileHash()); err != nil {
		return err
	}
	defer fw.close()
	if err := fw.createHeader(data.Header()); err != nil {
		return err
	}
	if err := fw.createDataFile(data); err != nil {
		os.Remove(fw.headerFile.Name())
		return err
	}
	return nil
}

func (fw *FileWriter) createHeader(header sfile.HeaderFormat) error {
	if b, err := json.Marshal(header.Attributes()); err == nil {
		if _, err = fw.headerFile.Write(b); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (fw *FileWriter) createDataFile(data sfile.DataFormat) error {
	log.Println("FileName to create:", fw.dataFile.Name())
	saveFile := bytes.NewBuffer([]byte(""))
	saveFile.WriteString("SAVE")
	// block write flags. 0 means not written, 1 means written
	// maybe need to pass back block count?
	blockCount := int(math.Ceil(float64(data.Size()) / float64(blockSize)))
	saveFile.Write(conversion.IntToBytes(blockCount))
	for i := 0; i < blockCount; i++ {
		saveFile.WriteByte('0')
	}
	saveFile.Write(conversion.Int64ToBytes(data.Size()))
	// Truncate file so that the file is created at the correct size.
	// This is beneficial when doing multiupload
	truncSize := int64(saveFile.Len()) + data.Size()
	if truncSize > maxFileSize {
		return errors.New("exceeded max file size allowed")
	}
	if err := fw.dataFile.Truncate(truncSize); err != nil {
		return err
	}
	if _, err := fw.dataFile.Write(saveFile.Bytes()); err != nil {
		os.Remove(fw.dataFile.Name())
		return err
	}
	return nil
}

// Update takes a SaveFile object and updates the contents of the header file and the data file.
func (fw *FileWriter) Update(data sfile.DataFormat) error {
	if err := fw.grabFilesForWrite(data.FileHash()); err != nil {
		return err
	}
	defer fw.close()
	if err := fw.updateHeader(data.Header()); err != nil {
		return err
	}
	if err := fw.updateDataFile(data); err != nil {
		return err
	}
	return nil
}

func (fw *FileWriter) updateHeader(header sfile.HeaderFormat) error {
	headerMap, err := fw.readHeader()
	if err != nil {
		return nil
	}
	for k, v := range header.Attributes() {
		headerMap[k] = v
	}
	newData, err := json.Marshal(headerMap)
	if err != nil {
		return err
	}
	_, err = fw.headerFile.Write(newData)
	return err
}

func (fw *FileWriter) updateDataFile(data sfile.DataFormat) error {
	if len(data.Data()) > blockSize {
		return errors.New("Data size is bigger than block allows")
	}
	fileData := make([]byte, 4)
	var offset int64 = 4
	// grab header size starting at position 4 skipping "SAVE" marker
	_, err := fw.dataFile.ReadAt(fileData, offset)
	if err != nil {
		return err
	}
	blockFlagSize := conversion.BytesToInt(fileData[0], fileData[1], fileData[2], fileData[3])
	if data.Block() > blockFlagSize {
		return errors.New("block out of range")
	}
	offset += 4
	blockOffset := offset
	blockFlags := make([]byte, blockFlagSize)
	_, err = fw.dataFile.ReadAt(blockFlags, offset)
	if err != nil {
		return err
	}
	offset += int64(blockFlagSize)
	// jump to block of data in file. 8 is for the data size bytes
	offset += 8 + int64(data.Block()*blockSize)
	_, err = fw.dataFile.WriteAt(data.Data(), offset)
	if err != nil {
		return err
	}
	blockFlags[data.Block()] = 1
	_, err = fw.dataFile.WriteAt(blockFlags, blockOffset)
	if err != nil {
		return err
	}
	return nil
}

// Delete SaveFile objects
func (fw *FileWriter) Delete(data sfile.DataFormat) error {
	pathBuffer := bytes.NewBufferString(fw.path)
	pathBuffer.WriteRune(os.PathSeparator)
	pathBuffer.Write(data.FileHash())
	if err := os.Remove(pathBuffer.String()); err != nil {
		return err
	}
	pathBuffer.WriteString(headerAppend)
	if err := os.Remove(pathBuffer.String()); err != nil {
		return err
	}
	return nil
}

// ReadHeader takes a SaveFile object and reads out the header file into the given SaveFile object
func (fw *FileWriter) ReadHeader(data sfile.DataFormat) error {
	if err := fw.grabFilesForWrite(data.FileHash()); err != nil {
		return err
	}
	defer fw.close()
	headerMap, err := fw.readHeader()
	if err != nil {
		return err
	}
	data.Header().SetAttributes(headerMap)
	return nil
}

func (fw *FileWriter) readHeader() (map[string]interface{}, error) {
	info, err := fw.headerFile.Stat()
	if err != nil {
		return nil, err
	}
	origData := make([]byte, info.Size())
	_, err = fw.headerFile.Read(origData)
	if err != nil {
		return nil, err
	}
	err = fw.headerFile.Truncate(info.Size())
	if err != nil {
		return nil, err
	}
	headerMap := make(map[string]interface{})
	err = json.Unmarshal(origData, headerMap)
	if err != nil {
		fw.headerFile.Write(origData)
		return nil, err
	}
	return headerMap, nil
}

// ReadDataBlock takes in a SaveFile object and reads the block of data specified.
// This implementation mutates the passed in SaveFile
func (fw *FileWriter) ReadDataBlock(data sfile.DataFormat) error {
	if err := fw.grabFilesForWrite(data.FileHash()); err != nil {
		return err
	}
	defer fw.close()
	if err := fw.readDataFile(data); err != nil {
		return err
	}
	return nil
}

func (fw *FileWriter) readDataFile(data sfile.DataFormat) error {
	fileData := make([]byte, 4)
	var offset int64 = 4
	// grab header size starting at position 4 skipping "SAVE" marker
	_, err := fw.dataFile.ReadAt(fileData, offset)
	if err != nil {
		return err
	}
	blockFlagSize := conversion.BytesToInt(fileData[0], fileData[1], fileData[2], fileData[3])
	if data.Block() > blockFlagSize {
		return errors.New("block out of range")
	}
	offset += 4 + int64(blockFlagSize)
	// jump to block of data in file. 8 is for the data size bytes
	offset += 8 + int64(data.Block()*blockSize)
	newData := make([]byte, blockSize)
	_, err = fw.dataFile.ReadAt(newData, offset)
	if err != nil {
		return err
	}
	return data.SetData(newData)
}

// close handles closing the FileWriter's internal objects.
func (fw *FileWriter) close() {
	err := fw.dataFile.Close()
	if err == nil {
		log.Printf("error closing file: %v\n", err)
	}
	err = fw.headerFile.Close()
	if err != nil {
		log.Printf("error closing file: %v\n", err)
	}
	fw.dataFile = nil
	fw.headerFile = nil
}

// FileExists check if SaveFile exists
func FileExists(fileName string) bool {
	exists := false
	if _, err := os.Stat(string(fileName)); !os.IsNotExist(err) {
		exists = true
	}
	if _, err := os.Stat(string(fileName + headerAppend)); !os.IsNotExist(err) {
		exists = true
	}
	return exists
}
