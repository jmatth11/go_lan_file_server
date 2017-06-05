package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sfile"
	"sort"
	"strings"
	"time"
)

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

func (g *GetFilesWithAttributes) sortedAttributeKeys() []string {
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

// CreateTodaysFolder is a method to create a folder with the current date as its name.
// @return string  The folder path
func CreateTodaysFolder() string {
	// Grab date
	year, month, day := time.Now().Date()
	// format string to desired file name
	name := fmt.Sprintf("Data\\%d-%d-%d", year, month, day)
	// check if folder already exists
	_, err := os.Stat(name)
	// if it doesn't exist, create it.
	if err != nil {
		os.Mkdir(name, 0666)
	}
	return name
}

// createHeaderObject creates a sfile.SimpleHeader object with default attributes
// @param data FileData object passed in to set the Attribute keys in the SimpleHeader object
// @return *sfile.SimpleHeader object
func createHeaderObject(data map[string]string) *sfile.SimpleHeader {
	headerObj := &sfile.SimpleHeader{Attributes: make(map[string]interface{})}
	for k, v := range data {
		headerObj.Attributes[k] = v
	}
	return headerObj
}

// PostFile is a method to handle post request for a file to be saved to the server.
func PostFile(w http.ResponseWriter, req *http.Request) {
	// create json decoder
	decoder := json.NewDecoder(req.Body)
	var data FileData
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Post File Error:", err)
	}
	defer req.Body.Close()
	headerObj := createHeaderObject(data.Attributes)
	filePath := bytes.NewBufferString(CreateTodaysFolder() + "\\")
	_, err = filePath.Write(data.ValidateFile)
	if err != nil {
		log.Print(err)
		errReturn := map[string]string{"Error": fmt.Sprintf("Error creating file path for %s; %s", data.ValidateFile, err)}
		writeOutJSONMessage(errReturn, w)
		return
	}
	n, err := sfile.WriteSaveFile(filePath.Bytes(), data.Data, headerObj, data.StartIndex, data.Size)
	if err != nil {
		log.Fatal(err)
		errReturn := map[string]string{"Error": fmt.Sprintf("Error while writing file %s; %s", data.ValidateFile, err)}
		writeOutJSONMessage(errReturn, w)
		return
	}
	log.Printf("File data, Name: %s. Wrote %d bytes", data.ValidateFile, n)
	writeOutJSONMessage(map[string]string{"Error": ""}, w)
}

// GetFolders is a method to retrieve the list of folder names in the Data path.
func GetFolders(w http.ResponseWriter, req *http.Request) {
	// Grab all folders in Data directory
	foldersFromDir, err := ioutil.ReadDir("Data\\.")
	if err != nil {
		log.Fatal(err)
		errFolders := FoldersList{Error: "ERROR: Could not read Data directory."}
		writeOutJSONMessage(errFolders, w)
		return
	}

	folders := FoldersList{Folders: make([]Folder, 0)}
	// Grab all folder info
	for _, obj := range foldersFromDir {
		// Grab files in folder
		subFiles, err := ioutil.ReadDir("Data\\" + obj.Name() + "\\.")
		if err != nil {
			log.Fatal(err)
		}
		// Create Folder object with its name and its file count
		f := Folder{Name: obj.Name(), Count: len(subFiles)}
		folders.Folders = append(folders.Folders, f)
	}
	writeOutJSONMessage(folders, w)
}

// GetFiles is a method to accept a POST request for a specific folder in the Data path requesting files
// from startIndex to endIndex(exclusively). Must also send a dictionary with the keys of the attributes you want to extract
func GetFiles(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data GetFilesWithAttributes
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Post File Error:", err)
	}
	defer req.Body.Close()
	// check if out of range
	if data.StartIndex > data.EndIndex || data.StartIndex < 0 || data.EndIndex < 0 {
		log.Fatal(fmt.Sprintf("ERROR: folder: %s, startIndex: %d, endIndex: %d", data.Folder, data.StartIndex, data.EndIndex))
		errFiles := FileDataList{Error: fmt.Sprintf("ERROR: Either start or end index is incorrect. startIndex: %d, endIndex: %d", data.StartIndex, data.EndIndex)}
		writeOutJSONMessage(errFiles, w)
		return
	}
	// get list of files from folder
	files, err := ioutil.ReadDir("Data\\" + data.Folder)
	if err != nil {
		log.Fatal(err)
		errFiles := FileDataList{Error: "ERROR: Folder given could not be opened. Folder: " + data.Folder}
		writeOutJSONMessage(errFiles, w)
		return
	}
	// create Header object with empty FileData object because it will be populated from the read
	headerObj := createHeaderObject(data.Attributes)
	allFiles := FileDataList{Files: make([]FileData, 0)}
	for _, obj := range files[data.StartIndex:data.EndIndex] {
		// create SaveFile object from reading file
		saveFileObj, err := sfile.ReadSaveFile([]byte(obj.Name()), headerObj)
		if err != nil {
			log.Fatal(err)
		}
		// base64 encode data
		dstData := make([]byte, base64.StdEncoding.EncodedLen(len(saveFileObj.Data)))
		base64.StdEncoding.Encode(dstData, saveFileObj.Data)
		// map our objects
		headerMap := saveFileObj.Header.GetHeader()
		for i, v := range data.sortedAttributeKeys() {
			data.Attributes[v] = headerMap[i]
		}
		// Split on path separator to grab file hash
		validFile := strings.Split(obj.Name(), "\\")
		// base64 encode file hash
		dstValid := make([]byte, base64.StdEncoding.EncodedLen(len(validFile[len(validFile)-1])))
		base64.StdEncoding.Encode(dstValid, []byte(validFile[len(validFile)-1]))
		// create our object
		f := FileData{Data: dstData, Size: int64(saveFileObj.Size), StartIndex: 0, ValidateFile: dstValid, Attributes: data.Attributes}
		allFiles.Files = append(allFiles.Files, f)
	}
	writeOutJSONMessage(allFiles, w)
}

/**
 * Method to take an object json.Marshal it and write it out
 * to the console and the reposewriter.
 * @param obj interface{} A struct value
 * @param w http.ResponseWriter
 */
func writeOutJSONMessage(obj interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}

func main() {
	// Check if Data Folder exists and if not, create it.
	_, err := os.Stat("Data")
	if err != nil {
		os.Mkdir("Data", 0666)
	}
	http.HandleFunc("/post_file", PostFile)
	http.HandleFunc("/get_folders", GetFolders)
	http.HandleFunc("/get_files", GetFiles)
	http.ListenAndServe(":8080", nil)
}
