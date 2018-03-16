package validation

import (
	"bytes"
	"crypto/sha256"
)

// FileValidator is a simple interface that enforces a validate method to compare data against a hash
type FileValidator interface {
	Validate(data, hash []byte) bool
}

// DefaultFileValidator is a simple object that utilizes 256 hash comparisons
type DefaultFileValidator struct {
	FileValidator
}

// Validate will perform a sha256 hash on the data and compare it against a sha256 hash given.
func (dv DefaultFileValidator) Validate(data, hash []byte) bool {
	checkHash := sha256.Sum256(data)
	if bytes.Compare(hash, checkHash[:]) != 0 {
		return false
	}
	return true
}
