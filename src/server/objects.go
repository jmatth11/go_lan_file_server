package server

// objects file to hold types used for the server

// Response is json wrapper for all objects to be sent back to user
type Response struct {
	Data    interface{}
	Error   string
	Success bool
}

// FileHeader holds information for updating file header info
type FileHeader struct {
	Attributes   map[string]string
	Path         string
	ValidateFile []byte
}

// FileInit should be File meta data sent to server/client
type FileInit struct {
	FileHeader
	BlockCount int
	Size       int64
}

// FileBlock represents single block of data sent to server/client
type FileBlock struct {
	Block        int
	Data         []byte
	Path         string
	ValidateFile []byte
}

// BadSaveFileResponse simple struct to represent a response for a user's request
type BadSaveFileResponse struct {
	Blocks []int
}

// Not sure if I want to change any of the objects below

// GetFilesWithAttributes is an object to hold the folder you wish to grab files from,
// the StartIndex and EndIndex of the files you want,
// And a map with the keys of the attributes you want to extract for the files
type GetFilesWithAttributes struct {
	Folder     string
	StartIndex int
	EndIndex   int
	Attributes map[string]string
}

// FoldersList is an object to store a list of Folder objects
type FoldersList struct {
	Folders []Folder
}

// Folder is an object to store the name of the folder and the count of files it holds
type Folder struct {
	Name  string
	Count int
}
