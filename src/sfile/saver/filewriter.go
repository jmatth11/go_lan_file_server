package saver

import (
	"bytes"
	"io"
	"os"
	"sfile"
)

const (
	headerAppend = "-Header"
)

// FileWriter object for saving a data file and its associated header attribute file.
// on a file system.
type FileWriter struct {
	SReaderWriter
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

func fileExists(filename string) bool {
	if _, err := os.Stat(string(fileName)); os.IsNotExist(err) {
		return false
	}
	return true
}

func saveDataFile(data *sfile.SaveFile, pos int64) (int64, error) {
	dataFile := string(data.FileHash)
	if fileExists(dataFile) {
		// TODO
	}
	return 0, nil
}

func saveHeaderFile(data *sfile.SaveFile) error {
	headerFile := string(data.FileHash) + headerAppend
	if fileExists(headerFile) {
		// TODO
	}
	return nil
}
