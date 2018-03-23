package saver

import "sfile"

// FileAccessor Interface defines basic methods to interact with a SaveFile
type FileAccessor interface {
	create(data *sfile.SaveFile) error

	update(data *sfile.SaveFile, pos int64) *sfile.SaveFileResponse

	delete(data *sfile.SaveFile) error

	read(path string) *sfile.SaveFileResponse
	// Write is used to save a SaveFile object.
	//Write(data *sfile.SaveFile, lastPos, pos int64) (int64, error)

	// Read
	//Read(data *sfile.SaveFile) error
}
