package sfile

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

// SimpleHeader is an object that just holds a map object of the attributes
// it needs to hold. This object also implements the HeaderFormat interface methods
type SimpleHeader struct {
	Attributes map[string]interface{}
}

// GetHeader is the method to grab the attributes out of the object
func (sh *SimpleHeader) GetHeader() map[string]interface{} {
	return sh.Attributes
}

// GetHeaderSize is the method to grab the size of the header for a byte slice
func (sh *SimpleHeader) GetHeaderSize() (n int, err error) {
	_, n, err = sh.bufferFromAttributes()
	return
}

// Read is the method that will go through the attributes in the Simpleheader object
// and populate the []byte parameter with a string representation of the attributes and their string sizes
func (sh *SimpleHeader) Read(b []byte) (n int, err error) {
	headBuf, n, err := sh.bufferFromAttributes()
	if err != nil {
		return
	}
	if cap(b) < n {
		errMsg := fmt.Sprintf("error: b only has capacity for %d while Header has size %d", cap(b), n)
		log.Fatal(errMsg)
		err = errors.New(errMsg)
		return
	}
	headBytes := headBuf.Bytes()
	for i := 0; i < cap(b); i++ {
		b[i] = headBytes[i]
	}
	return
}

// Write is the Method that extracts out the attributes and stores them as strings.
func (sh *SimpleHeader) Write(b []byte) (n int, err error) {
	if len(sh.Attributes) == 0 && len(b) > 0 {
		err = errors.New("error: Attributes field does not have any attributes to extract the header data")
		return
	}
	m := make(map[string]interface{})
	for k := range sh.Attributes {
		size := bytesToInt(b[n], b[n+1], b[n+2], b[n+3])
		n += 4
		m[k] = string(b[n : n+size])
		n += size
	}
	sh.Attributes = m
	return
}

func (sh *SimpleHeader) bufferFromAttributes() (headBuf *bytes.Buffer, n int, err error) {
	headBuf = bytes.NewBuffer([]byte(""))
	for _, v := range sh.Attributes {
		val := fmt.Sprintf("%s", v)
		count, errMsg := headBuf.Write(intToBytes(len(val)))
		if errMsg != nil {
			err = errMsg
			return
		}
		n += count
		count, errMsg = headBuf.WriteString(val)
		if errMsg != nil {
			err = errMsg
			return
		}
		n += count
	}
	return
}
