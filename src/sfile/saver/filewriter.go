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

// FileWriter object for saving a data file and its associated header attribute file.
// on a file system.
type FileWriter struct {
	dataFile   *os.File
	headerFile *os.File
}

// NewFileWriter creates a new FileWriter object and creates two files based on the
// parameter passed in. The main data file will be named the fileName passed in.
// The header attribute file will be named fileName + "-Header".
// @param fileName []byte The file name to use.
// @return (*FileWriter, error)
func NewFileWriter(fileName []byte) (*FileWriter, error) {
	f1, f2, err := grabFilesForWrite(fileName)
	if err != nil {
		return nil, err
	}
	return &FileWriter{
		dataFile:   f1,
		headerFile: f2,
	}, nil
}

// grabFilesForWrite creates the main data file and the associated header file.
func grabFilesForWrite(fileName []byte) (*os.File, *os.File, error) {
	nameBuffer := bytes.NewBuffer(fileName)
	file1, err := os.OpenFile(nameBuffer.String(), os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, nil, err
	}
	nameBuffer.WriteString(headerAppend)
	file2, err := os.OpenFile(nameBuffer.String(), os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, nil, err
	}
	return file1, file2, nil
}

// Create takes a SaveFile object and creates the header and data files needed
func (fw *FileWriter) Create(data *sfile.SaveFile) error {
	if err := fw.createHeader(data.Header); err != nil {
		return err
	}
	if err := fw.createDataFile(data); err != nil {
		os.Remove(fw.headerFile.Name())
		return err
	}
	return nil
}

func (fw *FileWriter) createHeader(header sfile.HeaderFormat) error {
	if b, err := json.Marshal(header); err == nil {
		if _, err = fw.headerFile.Write(b); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (fw *FileWriter) createDataFile(data *sfile.SaveFile) error {
	log.Println("FileName to create:", fw.dataFile.Name())
	saveFile := bytes.NewBuffer([]byte(""))
	saveFile.WriteString("SAVE")
	// block write flags. 0 means not written, 1 means written
	// maybe need to pass back block count?
	blockCount := int(math.Ceil(float64(data.Size) / float64(blockSize)))
	saveFile.Write(conversion.IntToBytes(blockCount))
	for i := 0; i < blockCount; i++ {
		saveFile.WriteByte('0')
	}
	saveFile.Write(conversion.Int64ToBytes(data.Size))
	// Truncate file so that the file is created at the correct size.
	// This is beneficial when doing multiupload
	truncSize := int64(saveFile.Len()) + data.Size
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
func (fw *FileWriter) Update(data *sfile.SaveFile) error {
	if err := fw.updateHeader(data.Header); err != nil {
		return err
	}
	if err := fw.updateDataFile(data); err != nil {
		return err
	}
	return nil
}

func (fw *FileWriter) updateHeader(header sfile.HeaderFormat) error {
	fw.headerFile.Seek(0, 0)
	info, err := fw.headerFile.Stat()
	if err != nil {
		return err
	}
	origData := make([]byte, info.Size())
	_, err = fw.headerFile.Read(origData)
	if err != nil {
		return err
	}
	err = fw.headerFile.Truncate(info.Size())
	if err != nil {
		return err
	}
	headerMap := make(map[string]interface{})
	err = json.Unmarshal(origData, headerMap)
	if err != nil {
		fw.headerFile.Write(origData)
		return err
	}
	for k, v := range header.Attributes() {
		headerMap[k] = v
	}
	newData, err := json.Marshal(headerMap)
	if err != nil {
		fw.headerFile.Write(origData)
		return err
	}
	_, err = fw.headerFile.Write(newData)
	return err
}

func (fw *FileWriter) updateDataFile(data *sfile.SaveFile) error {
	if len(data.Data) > blockSize {
		return errors.New("Data size is bigger than block allows")
	}
	fw.dataFile.Seek(0, 0)
	fileData := make([]byte, 4)
	var offset int64 = 4
	// grab header size starting at position 4 skipping "SAVE" marker
	_, err := fw.dataFile.ReadAt(fileData, offset)
	if err != nil {
		return err
	}
	blockFlagSize := conversion.BytesToInt(fileData[0], fileData[1], fileData[2], fileData[3])
	if data.Block > blockFlagSize {
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
	offset += 8 + int64(data.Block*blockSize)
	_, err = fw.dataFile.WriteAt(data.Data, offset)
	if err != nil {
		return err
	}
	blockFlags[data.Block] = 1
	_, err = fw.dataFile.WriteAt(blockFlags, blockOffset)
	if err != nil {
		return err
	}
	return nil
}

// Delete SaveFile objects
func (fw *FileWriter) Delete(data *sfile.SaveFile) error {
	// TODO
	return nil
}

// Read SaveFile object at path
func (fw *FileWriter) Read(path string) *sfile.SaveFile {
	// TODO
	return nil
}

// Close handles closing the FileWriter's internal objects.
// @return error
func (fw *FileWriter) Close() error {
	err := fw.dataFile.Close()
	if err != nil {
		return err
	}
	err = fw.headerFile.Close()
	return err
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
