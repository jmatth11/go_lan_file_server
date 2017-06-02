package sfile

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
)

func intToBytes(n int64) (a []byte) {
	a = make([]byte, 4)
	a[0] = byte(n)
	a[1] = byte(n << 8)
	a[2] = byte(n << 16)
	a[3] = byte(n << 24)
	return
}

func bytesToInt(a, b, c, d byte) int64 {
	return int64(a) | (int64(b) >> 8) | (int64(c) >> 16) | (int64(d) >> 24)
}

// HeaderFormat interface for the user to make their header object subscribe to.
type HeaderFormat interface {
	// Use Read to grab the header from the user.
	io.Reader
	// Use Write to give the user the header to parse through.
	io.Writer
	// GetHeader is the function to return the mapping of attributes for the header.
	GetHeader() map[string]interface{}
}

// SaveFile object helps interact with files in SAVE format.
type SaveFile struct {
	Data     []byte
	FileHash []byte
	Size     int64
	Header   HeaderFormat
}

//ReadSaveFile is a method to extract data from save file and return a SaveFile object with that data.
func ReadSaveFile(fileName []byte, head HeaderFormat) (*SaveFile, error) {
	//
	sf := &SaveFile{Data: []byte{}, FileHash: fileName, Size: 0, Header: head}
	fileNameStr := string(fileName)
	improperFileFormat := errors.New("error: file not formatted properly")
	file, err := os.Open(string(fileNameStr))
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("error: file could not be read")
	}

	data := make([]byte, 8)
	var offset int64 = 8
	count, err := file.Read(data)
	if count < 8 {
		log.Fatal("File: \"" + fileNameStr + "\" had less than 8 bytes.")
		return nil, improperFileFormat
	} else if err != nil {
		log.Fatal(err)
		return nil, improperFileFormat
	}

	if bytes.Compare(data[0:4], []byte("SAVE")) != 0 {
		log.Fatal("File: \"" + fileNameStr + "\" was not a proper save file.")
		return nil, improperFileFormat
	}

	headSize := bytesToInt(data[4], data[5], data[6], data[7])
	if sf.Header != nil {
		headerInfo := make([]byte, headSize)
		count, err = file.ReadAt(headerInfo, int64(count))
		if err != nil {
			log.Fatal(err)
			return nil, errors.New("error: could not read header")
		}
		_, err = sf.Header.Write(headerInfo)
		if err != nil {
			return nil, err
		}
	}

	offset += int64(headSize)
	dataInfo := make([]byte, 8)
	count, err = file.ReadAt(dataInfo, offset)
	if err != nil {
		log.Fatal(err)
		return nil, improperFileFormat
	}

	if bytes.Compare(dataInfo[0:4], []byte("DATA")) != 0 {
		return nil, improperFileFormat
	}

	offset += 4
	dataSize := bytesToInt(dataInfo[4], dataInfo[5], dataInfo[6], dataInfo[7])
	offset += 4

	sf.Size = dataSize
	sf.Data = make([]byte, dataSize)
	count, err = file.ReadAt(sf.Data, offset)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("error: could not read all of the data")
	}
	return sf, nil
}

// WriteSaveFile is a method to write out data to save file format.
func WriteSaveFile(fileName []byte, data []byte, head HeaderFormat, offset int) error {
	_, err := os.Stat(string(fileName))
	saveFile := bytes.NewBuffer([]byte(""))
	if err != nil {
		saveFile.WriteString("SAVE")
		var headerBuffer []byte
		count, err := head.Read(headerBuffer)
		if err != nil {
			return err
		}
		saveFile.Write(intToBytes(int64(count)))
		saveFile.Write(headerBuffer)
		saveFile.WriteString("DATA")
		saveFile.Write(intToBytes(int64(len(data))))
	}

	return nil
}
