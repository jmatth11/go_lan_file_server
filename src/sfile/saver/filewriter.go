package saver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sfile"
)

const (
	headerAppend = "-Header"
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

func (fw *FileWriter) Write(data *sfile.SaveFile, pos int64) (newPos int64, err error) {
	newPos, err = saveDataFile(data, pos)
	if err != nil {
		return
	}
	err = saveHeaderFile(data)
	if err != nil {
		return
	}
	return
}

func (fw *FileWriter) Read(data *sfile.SaveFile) error {
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

func fileExists(fileName string) bool {
	if _, err := os.Stat(string(fileName)); os.IsNotExist(err) {
		return false
	}
	return true
}

func saveDataFile(data *sfile.SaveFile, lastPos, pos int64) (int64, error) {
	dataFile := string(data.FileHash)
	fileObj, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE, 0777)
	newPos := 0
	if err != nil {
		return 0, err
	}
	defer fileObj.Close()
	if fileExists(dataFile) {
		fileData := make([]byte, 4)
		offset := 4
		// grab header size starting at position 4 skipping "SAVE" marker
		_, err = fileObj.ReadAt(fileData, int64(offset))
		if err != nil {
			return 0, err
		}
		origSize := sfile.BytesToInt(fileData[0], fileData[1], fileData[2], fileData[3])
		if origSize == data.Size {
			errMsg := fmt.Sprintf("error: the size of the data matches the size of the original file. The Entire file should already exist.")
			return 0, errors.New(errMsg)
		}
		// TODO don't know if I want to allow out of order block writes yet.
		if origSize != lastPos {
			errMsg := fmt.Sprintf("error: last received index was %d; current received index was %d", origSize, lastPos)
			return origSize, errors.New(errMsg)
		}
		newDataSize := origSize + len(data.Data)
		newPos = newDataSize
		_, err = fileObj.WriteAt(intToBytes(newDataSize), int64(offset))
		offset += 4
		if err != nil {
			return 0, err
		}
		_, err = fileObj.WriteAt(data.Data, int64(offset+origSize))
		if err != nil {
			return 0, err
		}
	} else {
		log.Println("FileName to create:", dataFile)
		saveFile := bytes.NewBuffer([]byte(""))
		saveFile.WriteString("SAVE")
		dataSize := len(data.Data)
		newPos = dataSize
		saveFile.Write(intToBytes(dataSize))
		// Truncate file so that the file is created at the correct size.
		// This is beneficial when doing multiupload
		err = fileObj.Truncate(int64(saveFile.Len()) + data.Size)
		if err != nil {
			return 0, err
		}
		saveFile.Write(data.Data)
		_, err = fileObj.Write(saveFile.Bytes())
		if err != nil {
			return 0, err
		}
	}
	return newPos, nil
}

func saveHeaderFile(data *sfile.SaveFile) error {
	headerFile := string(data.FileHash) + headerAppend
	if fileExists(headerFile) {
		b, err := ioutil.ReadFile(headerFile)
		if err != nil {
			return err
		}
		origData := make(map[string]interface{})
		err = json.Unmarshal(b, origData)
		if err != nil {
			return err
		}
		for k, v := range data.Header {
			origData[k] = v
		}
		b, err = json.Marshal(origData)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(headerFile, b, os.ModePerm); err != nil {
			return err
		}
	} else {
		if b, err := json.Marshal(data.Header); err == nil {
			if err = ioutil.WriteFile(headerFile, b, os.ModePerm); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
