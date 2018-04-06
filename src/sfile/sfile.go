package sfile

// DataFormat is the interface that Save File implementations should conform to.
type DataFormat interface {
	Data() []byte
	SetData(data []byte) error

	FileHash() []byte
	SetFileHash(fileHash []byte) error

	Size() int64
	SetSize(size int64)

	Header() HeaderFormat
	SetHeader(header HeaderFormat) error

	Block() int
	SetBlock(block int)
}

// SimpleSaveFile object helps interact with files in SAVE format.
type SimpleSaveFile struct {
	// Actual File Data
	data []byte
	// The Hash of the Data of the file, which is also the name of the file on the server
	fileHash []byte
	// The Size of the Data
	size int64
	// The Header for the File. Contains Attributes of the file (name, type, etc..)
	header HeaderFormat
	// The block of data to write/read
	block int
}

// NewSimpleSaveFile returns a newly generated SimpleSaveFile object
func NewSimpleSaveFile() *SimpleSaveFile {
	return new(SimpleSaveFile)
}

// Data returns the SimpleSaveFile's data
func (ssf *SimpleSaveFile) Data() []byte {
	return ssf.data
}

// SetData sets the SimpleSaveFile's data
func (ssf *SimpleSaveFile) SetData(data []byte) error {
	ssf.data = data
	return nil
}

// FileHash returns the SimpleSaveFile's file hash.
// This is also used as the file name.
func (ssf *SimpleSaveFile) FileHash() []byte {
	return ssf.fileHash
}

// SetFileHash sets the SimpleSaveFile's fileHash.
// This is also used as the file name.
func (ssf *SimpleSaveFile) SetFileHash(fileHash []byte) error {
	ssf.fileHash = fileHash
	return nil
}

// Size returns the size of the file
func (ssf *SimpleSaveFile) Size() int64 {
	return ssf.size
}

// SetSize sets the Size for the file
func (ssf *SimpleSaveFile) SetSize(size int64) {
	ssf.size = size
}

// Header returns the SimpleSaveFile's header format object
func (ssf *SimpleSaveFile) Header() HeaderFormat {
	return ssf.header
}

// SetHeader sets the SimpleSaveFile's header format object
func (ssf *SimpleSaveFile) SetHeader(header HeaderFormat) error {
	ssf.header = header
	return nil
}

// Block returns the SimpleSaveFile's current block
func (ssf *SimpleSaveFile) Block() int {
	return ssf.block
}

// SetBlock sets the SimpleSaveFile's current block
func (ssf *SimpleSaveFile) SetBlock(block int) {
	ssf.block = block
}
