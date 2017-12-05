package saver

import "sfile"

// SReaderWriter interface contains methods to Write and Read SaveFile objects.
type SReaderWriter interface {
	// Write is used to save a SaveFile object.
	Write(data *sfile.SaveFile, pos int64) (int64, error)

	// Read
	Read(data *sfile.SaveFile) error
}
