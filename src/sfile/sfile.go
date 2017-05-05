package sfile

import (
  "log"
  "os"
  "errors"
)

func intToBytes(uint64 n) (a, b, c, d byte) {
  a = byte(n)
  b = byte(n << 8)
  c = byte(n << 16)
  d = byte(n << 24)
  return
}

func bytesToInt(a, b, c, d byte) uint64 {
  return uint64(a) | uint64(b >> 8) | uint64(c >> 16) | uint64(d >> 24)
}

type HeaderFormat interface {
  ParseHeader(file *os.File, data []byte) error
}

type SaveFile struct {
	Data      []byte
	FileHash  []byte
	Size 		  uint64
	Header    HeaderFormat
}

func ReadSaveFile(file_name string, head HeaderFormat) (*SaveFile, error) {
  sf := &SaveFile{Data:make([]byte, 0), FileHash:make([]byte, 0), Size:0, Header:head}
  improper_file_format := errors.New("ERROR: File not formatted properly.")
  file, err := os.Open(file_name)
  if err != nil {
    log.Fatal(err)
    return nil, erros.New("ERROR: File could not be read.")
  }
  data := make([]byte, 8)
  count, err := file.Read(data)
  if count < 8 {
    log.Fatal("File: \"" + file_name + "\" had less than 8 bytes.")
    return nil, improper_file_format
  } else if err != nil {
    log.Fatal(err)
    return nil, improper_file_format
  }
  if data[0:4] != []byte("SAVE") {
    log.Fatal("File: \"" + file_name + "\" was not a proper save file.")
    return nil, improper_file_format
  }
  head_size := bytesToInt(data[4:8]...)
  if sf.Header != nil {
    header_info := make([]byte, head_size)
    count, err = file.ReadAt(header_info, count)
    if err != nil {
      log.Fatal(err)
      return nil, errors.New("ERROR: Could not read header.")
    }
    err = sf.Header.ParseHeader(file, header_info)
    if err != nil {
      return nil, err
    }
  }
  
}
