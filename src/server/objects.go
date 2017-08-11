package server

// objects file to hold types used for the server

import "sort"

// FileData is an object that represents all the data we store for a file saved
type FileData struct {
	Data         []byte
	ValidateFile []byte
	StartIndex   int
	Size         int64
	Attributes   map[string]string
}

// GetFilesWithAttributes is an object to hold the folder you wish to grab files from,
// the StartIndex and EndIndex of the files you want,
// And a map with the keys of the attributes you want to extract for the files
type GetFilesWithAttributes struct {
	Folder     string
	StartIndex int
	EndIndex   int
	Attributes map[string]string
}

// SortedAttributeKeys method sorts attributes
func (g *GetFilesWithAttributes) SortedAttributeKeys() []string {
	sortedKeys := make([]string, len(g.Attributes))
	i := 0
	for k := range g.Attributes {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

// FileDataList is an object to store a list of FileData objects
type FileDataList struct {
	Files []FileData
	Error string
}

// FoldersList is an object to store a list of Folder objects
type FoldersList struct {
	Folders []Folder
	Error   string
}

// Folder is an object to store the name of the folder and the count of files it holds
type Folder struct {
	Name  string
	Count int
}
