package sfile

import (
  "fmt"
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
  ParseHeader(data []byte)
}

type SaveFile struct {
	Data      []byte
	FileHash  []byte
	Size 		  int64
	Header    HeaderFormat
}

func ReadSaveFile(file_name string) *SaveFile {

}
