package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"./src/sfile"
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
		log.Printf("Created folder %s", name)
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
	logServerCall(req, "PostFile")
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
		errReturn := map[string]interface{}{"Count": 0, "Error": fmt.Sprintf("Error creating file path for %s; %s", data.ValidateFile, err)}
		writeOutJSONMessage(errReturn, w)
		return
	}
	n, err := sfile.WriteSaveFile(filePath.Bytes(), data.Data, headerObj, data.StartIndex, data.Size)
	if err != nil {
		log.Fatal(err)
		errReturn := map[string]interface{}{"Count": n, "Error": fmt.Sprintf("Error while writing file %s; %s", data.ValidateFile, err)}
		writeOutJSONMessage(errReturn, w)
		return
	}
	log.Printf("File data, Name: %s. Wrote %d bytes", data.ValidateFile, n)
	writeOutJSONMessage(map[string]interface{}{"Count": n, "Error": ""}, w)
}

// GetFolders is a method to retrieve the list of folder names in the Data path.
func GetFolders(w http.ResponseWriter, req *http.Request) {
	logServerCall(req, "GetFolders")
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
	logServerCall(req, "GetFiles")
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

// ValidateFile is a GET request that takes in a file hash and checks to see
// if that file exists on the server as a whole.
func ValidateFile(w http.ResponseWriter, req *http.Request) {
	logServerCall(req, "ValidateFile")
	errMsg := map[string]interface{}{"Error": "", "Size": 0}
	folder := req.URL.Query().Get("folder")
	fileHash := req.URL.Query().Get("hash")
	if fileHash == "" {
		index, err := strconv.Atoi(req.URL.Query().Get("index"))
		if err != nil {
			errMsg["Error"] = err
			writeOutJSONMessage(errMsg, w)
			return
		}
		validateFileWithIndex(w, req, folder, index)
	} else {
		validateFileWithHash(w, req, folder, fileHash)
	}
}

func validateFileWithIndex(w http.ResponseWriter, req *http.Request, folder string, index int) {
	log.Printf("validating %d from %s", index, folder)
	errMsg := map[string]interface{}{"Error": "", "Size": 0}
	files, err := ioutil.ReadDir("Data\\" + folder)
	if err != nil {
		errMsg["Error"] = err
		writeOutJSONMessage(errMsg, w)
		return
	}
	if index >= len(files) || index < 0 {
		errMsg["Error"] = errors.New("error: index out of range")
		writeOutJSONMessage(errMsg, w)
		return
	}
	saveFileObj, err := sfile.ReadSaveFile([]byte(files[index].Name()), nil)
	if err != nil {
		errMsg["Error"] = err
		writeOutJSONMessage(errMsg, w)
		return
	}
	splitFilePath := strings.Split(files[index].Name(), "\\")
	correctHash := []byte(splitFilePath[len(splitFilePath)-1])
	checkHash := sha256.Sum256(saveFileObj.Data)
	if bytes.Compare(correctHash, checkHash[:]) == 0 {
		writeOutJSONMessage(errMsg, w)
		return
	}
	errMsg["Error"] = errors.New("error: original hash does not match current data hash")
	errMsg["Size"] = saveFileObj.Size
	writeOutJSONMessage(errMsg, w)
}

func validateFileWithHash(w http.ResponseWriter, req *http.Request, folder, hash string) {
	log.Printf("validating %s from %s", hash, folder)
	errMsg := map[string]interface{}{"Error": "", "Size": 0}
	saveFileObj, err := sfile.ReadSaveFile([]byte("Data\\"+folder+hash), nil)
	if err != nil {
		errMsg["Error"] = err
		writeOutJSONMessage(errMsg, w)
		return
	}
	correctHash := []byte(hash)
	checkHash := sha256.Sum256(saveFileObj.Data)
	if bytes.Compare(correctHash, checkHash[:]) == 0 {
		writeOutJSONMessage(errMsg, w)
		return
	}
	errMsg["Error"] = errors.New("error: original hash does not match current data hash")
	errMsg["Size"] = saveFileObj.Size
	writeOutJSONMessage(errMsg, w)
}

// PingServ method listens for any message and sends back a response that lets
// the user know it is hitting the right address.
func PingServ(w http.ResponseWriter, req *http.Request) {
	logServerCall(req, "PingServ")
	w.Write([]byte("connected"))
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
	log.Printf("writeOutJSONMessage: %s", string(b))
	w.Write(b)
}

func logServerCall(req *http.Request, funcName string) {
	// right now just logging direct ip.
	// later might want to add req.Header.Get("X-Forwarded-For") to get possible tail of ips.
	log.Printf("%s|%s|%s %s", req.Method, funcName, "directly from:", req.RemoteAddr)
}

func main() {
	// Check if Data Folder exists and if not, create it.
	log.Println("starting up server on port :8080")
	_, err := os.Stat("Data")
	if err != nil {
		log.Println("creating Initial Data folder")
		os.Mkdir("Data", 0666)
	}
	http.HandleFunc("/ping", PingServ)
	http.HandleFunc("/post_file", PostFile)
	http.HandleFunc("/get_folders", GetFolders)
	http.HandleFunc("/get_files", GetFiles)
	http.HandleFunc("/validate_file", ValidateFile)
	http.ListenAndServe(":8080", nil)
}
