package sfile

import (
	"io"
)

func IntToBytes(n int) (a []byte) {
	a = make([]byte, 4)
	a[0] = byte(n)
	a[1] = byte(n >> 8)
	a[2] = byte(n >> 16)
	a[3] = byte(n >> 24)
	return
}

func BytesToInt(a, b, c, d byte) int {
	return int(a) | (int(b) << 8) | (int(c) << 16) | (int(d) << 24)
}

// HeaderFormat interface for the user to make their header object subscribe to.
type HeaderFormat interface {
	// Use Read to grab the header from the user.
	io.Reader
	// Use Write to give the user the header to parse through.
	io.Writer
	// GetHeader is the function to return the mapping of attributes for the header.
	// the interface value needs to implement a String method to be stringified.
	GetHeader() []string
	// GetHeaderSize is the function to return the size of the header for byte slices
	GetHeaderSize() (n int, err error)
}

// SaveFile object helps interact with files in SAVE format.
type SaveFile struct {
	// Actual File Data
	Data []byte
	// The Hash of the Data of the file, which is also the name of the file on the server
	FileHash []byte
	// The Size of the Data
	Size int
	// The Header for the File. Contains Attributes of the file (name, type, etc..)
	Header map[string]interface{}
	// The block of data to write/read
	Block int
}

// ReadSaveFile is a method to extract data from save file and return a SaveFile object with that data.
// Method should only be used when user is querying for file to return a SaveFile object
// that can be sent as json to user.
// func ReadSaveFile(fileName []byte, head HeaderFormat) (*SaveFile, error) {
// 	log.Printf("accessing file for read: %s", string(fileName))
// 	//!!! broken now because I change the struct signature !!!!
// 	sf := &SaveFile{Data: []byte{}, FileHash: fileName, Size: 0, Header: head}
// 	fileNameStr := string(fileName)
// 	improperFileFormat := errors.New("error: file not formatted properly")
// 	file, err := os.Open(string(fileNameStr))
// 	if err != nil {
// 		return nil, errors.New("error: file could not be read")
// 	}
// 	defer file.Close()
// 	data := make([]byte, 8)
// 	var offset int64 = 8
// 	count, err := file.Read(data)
// 	if count < 8 {
// 		return nil, improperFileFormat
// 	} else if err != nil {
// 		return nil, improperFileFormat
// 	}

// 	if bytes.Compare(data[0:4], []byte("SAVE")) != 0 {
// 		return nil, improperFileFormat
// 	}

// 	headSize := bytesToInt(data[4], data[5], data[6], data[7])
// 	if sf.Header != nil {
// 		headerInfo := make([]byte, headSize)
// 		count, err = file.ReadAt(headerInfo, int64(count))
// 		if err != nil {
// 			return nil, errors.New("error: could not read header")
// 		}
// 		_, err = sf.Header.Write(headerInfo)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	offset += int64(headSize)
// 	dataInfo := make([]byte, 8)
// 	count, err = file.ReadAt(dataInfo, offset)
// 	if err != nil {
// 		return nil, improperFileFormat
// 	}

// 	if bytes.Compare(dataInfo[0:4], []byte("DATA")) != 0 {
// 		return nil, improperFileFormat
// 	}

// 	offset += 4
// 	dataSize := bytesToInt(dataInfo[4], dataInfo[5], dataInfo[6], dataInfo[7])
// 	offset += 4

// 	sf.Size = dataSize
// 	sf.Data = make([]byte, dataSize)
// 	count, err = file.ReadAt(sf.Data, offset)
// 	if err != nil {
// 		return nil, errors.New("error: could not read all of the data")
// 	}
// 	return sf, nil
// }

// // WriteSaveFile is a method to write out data to save file format.
// // This method should only be used to take data from user and write to file.
// func WriteSaveFile(fileName []byte, data []byte, head HeaderFormat, lastPos int, size int64) (int, error) {
// 	log.Printf("accessing file for write: %s", string(fileName))
// 	_, fileAlreadyExists := os.Stat(string(fileName))
// 	fileObj, err := os.OpenFile(string(fileName), os.O_RDWR|os.O_CREATE, 0777)
// 	newPos := 0
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer fileObj.Close()
// 	if fileAlreadyExists != nil {
// 		log.Println("FileName to create:", string(fileName))
// 		saveFile := bytes.NewBuffer([]byte(""))
// 		saveFile.WriteString("SAVE")
// 		headerSize, err := head.GetHeaderSize()
// 		if err != nil {
// 			return 0, err
// 		}
// 		headerBuffer := make([]byte, headerSize)
// 		count, err := head.Read(headerBuffer)
// 		if err != nil {
// 			return 0, err
// 		}
// 		saveFile.Write(intToBytes(count))
// 		saveFile.Write(headerBuffer)
// 		saveFile.WriteString("DATA")
// 		dataSize := len(data)
// 		newPos = dataSize
// 		saveFile.Write(intToBytes(dataSize))
// 		// Truncate file so that the file is created at the correct size.
// 		// This is beneficial when doing multiupload
// 		err = fileObj.Truncate(int64(saveFile.Len()) + size)
// 		if err != nil {
// 			return 0, err
// 		}
// 		saveFile.Write(data)
// 		_, err = fileObj.Write(saveFile.Bytes())
// 		if err != nil {
// 			return 0, err
// 		}
// 	} else {
// 		fileData := make([]byte, 4)
// 		// grab header size starting at position 4 skipping "SAVE" marker
// 		_, err = fileObj.ReadAt(fileData, 4)
// 		if err != nil {
// 			return 0, err
// 		}
// 		offset := bytesToInt(fileData[0], fileData[1], fileData[2], fileData[3])
// 		// add 12 to offset to data size. this accounts for "SAVE", header size, and "DATA"
// 		offset += 12
// 		fileData = make([]byte, 4)
// 		// grab DATA size
// 		_, err = fileObj.ReadAt(fileData, int64(offset))
// 		origSize := bytesToInt(fileData[0], fileData[1], fileData[2], fileData[3])
// 		if int64(origSize) == size {
// 			errMsg := fmt.Sprintf("error: the size of the data matches the size of the original file. The Entire file should already exist.")
// 			return 0, errors.New(errMsg)
// 		}
// 		if origSize != lastPos {
// 			errMsg := fmt.Sprintf("error: last received index was %d; current received index was %d", origSize, lastPos)
// 			return origSize, errors.New(errMsg)
// 		}
// 		newDataSize := origSize + len(data)
// 		newPos = newDataSize
// 		_, err = fileObj.WriteAt(intToBytes(newDataSize), int64(offset))
// 		offset += 4
// 		if err != nil {
// 			return 0, err
// 		}
// 		_, err = fileObj.WriteAt(data, int64(offset+origSize))
// 		if err != nil {
// 			return 0, err
// 		}
// 	}
// 	return newPos, nil
// }
