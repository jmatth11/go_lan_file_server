package saver

import "sfile"

// FileAccessor Interface defines basic methods to interact with a SaveFile
type FileAccessor interface {
	Create(data *sfile.SaveFile) error

	Update(data *sfile.SaveFile, pos int64) *sfile.SaveFileResponse

	Delete(data *sfile.SaveFile) error

	ReadHeader(data *sfile.SaveFile) error

	ReadDataBlock(data *sfile.SaveFile) error
}
