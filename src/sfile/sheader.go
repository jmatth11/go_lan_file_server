package sfile

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sort"
)

// SimpleHeader is an object that just holds a map object of the attributes
// it needs to hold and implements the HeaderFormat interface.
// The Attributes are saved in the map's keys' alphabetical order
type SimpleHeader struct {
	Attributes map[string]interface{}
}

// GetHeader is the method to grab the attributes out of the object.
// The values of the Attribute's map are returned in the keys' alphabetical order.
func (sh *SimpleHeader) GetHeader() []string {
	headerList := make([]string, len(sh.Attributes))
	for i, v := range sh.sortedAttributeKeys() {
		headerList[i] = fmt.Sprintf("%s", sh.Attributes[v])
	}
	return headerList
}

// GetHeaderSize is the method to grab the size of the header for a byte slice
func (sh *SimpleHeader) GetHeaderSize() (n int, err error) {
	_, n, err = sh.bufferFromAttributes()
	return
}

// Read is the method that will go through the attributes in the Simpleheader object
// and populate the []byte parameter with a string representation of the attributes and their string sizes
// Read SimpleHeader description to see how attributes are written to.
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
// Read SimpleHeader description to see how attributes are written to.
func (sh *SimpleHeader) Write(b []byte) (n int, err error) {
	bLength := len(b)
	if len(sh.Attributes) == 0 && bLength > 0 {
		err = errors.New("error: Attributes field does not have any attributes to extract the header data")
		return
	}
	for _, k := range sh.sortedAttributeKeys() {
		if n >= bLength {
			return
		}
		size := bytesToInt(b[n], b[n+1], b[n+2], b[n+3])
		n += 4
		sh.Attributes[k] = string(b[n : n+size])
		n += size
	}
	return
}

func (sh *SimpleHeader) bufferFromAttributes() (headBuf *bytes.Buffer, n int, err error) {
	headBuf = bytes.NewBuffer([]byte(""))
	for _, v := range sh.sortedAttributeKeys() {
		val := fmt.Sprintf("%s", sh.Attributes[v])
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

func (sh *SimpleHeader) sortedAttributeKeys() []string {
	sortedKeys := make([]string, len(sh.Attributes))
	i := 0
	for k := range sh.Attributes {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}
