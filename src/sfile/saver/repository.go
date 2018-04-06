package saver

import "sfile"

// FileAccessor Interface defines basic methods to interact with a SaveFile
type FileAccessor interface {
	Create(data *sfile.DataFormat) error

	Update(data *sfile.DataFormat, pos int64) error

	Delete(data *sfile.DataFormat) error

	ReadHeader(data *sfile.DataFormat) error

	ReadDataBlock(data *sfile.DataFormat) error
}
