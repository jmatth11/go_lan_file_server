package sfile

import (
	"bytes"
	"io"
	"os"
)

const (
	headerAppend = "-Header"
)

// SWriter interface contains methods to save data and header files.
type SWriter interface {
	// SaveData is used to save a block of data from a certain position.
	SaveData(data []byte, pos int64) (int64, error)

	// SaveHeader is used to save header attributes.
	SaveHeader(data map[string]interface{}) error
}

// FileWriter object for saving a data file and its associated header attribute file.
// on a file system.
type FileWriter struct {
	dataFile   *os.File // main data file
	headerFile *os.File // file for header attributes
	io.Closer           // conform to io.Closer interface
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
