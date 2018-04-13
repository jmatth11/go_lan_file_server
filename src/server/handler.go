package server

// handler file to hold the server logic

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
	"path/filepath"
	"server/logger"
	"sfile"
	"strconv"
	"strings"
)

var rootPath = "Data"

// Practically everything needs to change in this file.
// might can reuse some code though...

func New(rootDir string) http.Handler {
	mux := http.NewServeMux()
	if strings.TrimSpace(rootDir) != "" {
		rootPath = rootDir
	}
	_, err := os.Stat(rootPath)
	if err != nil {
		Logln("creating Initial Data folder")
		os.Mkdir(rootPath, 0777)
	}
	// register functions
	mux.HandleFunc("/ping", pingServ)
	mux.HandleFunc("/post_file", postFile)
	mux.HandleFunc("/get_folders", getFolders)
	mux.HandleFunc("/get_files", getFiles)
	mux.HandleFunc("/validate_file", validateFile)
	return mux
}

// PostFile is a method to handle post request for a file to be saved to the server.
func postFile(w http.ResponseWriter, req *http.Request) {
	WriteFile(w, req)
}

// GetFolders is a method to retrieve the list of folder names in the Data path.
func getFolders(w http.ResponseWriter, req *http.Request) {
	GetFolders(w, req)
}

// GetFiles is a method to accept a POST request for a specific folder in the Data path requesting files
// from startIndex to endIndex(exclusively). Must also send a dictionary with the keys of the attributes you want to extract
func getFiles(w http.ResponseWriter, req *http.Request) {
	GetFiles(w, req)
}

// ValidateFile is a GET request that takes in a file hash and checks to see
// if that file exists on the server as a whole.
func validateFile(w http.ResponseWriter, req *http.Request) {
	ValidateFile(w, req)
}

// PingServ method listens for any message and sends back a response that lets
// the user know it is hitting the right address.
func pingServ(w http.ResponseWriter, req *http.Request) {
	PingServ(w, req)
}

// CreateTodaysFolder is a method to create a folder with the current date as its name.
// @return string  The folder path
// func CreateTodaysFolder() string {
// 	// Grab date
// 	year, month, day := time.Now().Date()
// 	// format string to desired file name
// 	name := fmt.Sprintf(filepath.Join(rootPath, "%d-%d-%d"), year, month, day)
// 	// check if folder already exists
// 	_, err := os.Stat(name)
// 	// if it doesn't exist, create it.
// 	if err != nil {
// 		//Logf("Created folder %s", name)
// 		os.Mkdir(name, 0777)
// 	}
// 	return name
// }

// PingServ method listens for any message and sends back a response that lets
// the user know it is hitting the right address.
func PingServ(w http.ResponseWriter, req *http.Request) {
	logger.Trace.ServerCall(req, "PingServ")
	w.Write([]byte("connected"))
}

// WriteFile is a method to handle post request for a file to be saved to the
// func WriteFile(w http.ResponseWriter, req *http.Request) {
// 	if req.Method != http.MethodPost {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		response := Response{Error: "Only POST method allowed."}
// 		WriteOutJSONMessage(response, w)
// 		return
// 	}
// 	LogServerCall(req, "WriteFile")
// 	// create json decoder
// 	decoder := json.NewDecoder(req.Body)
// 	var data FileData
// 	err := decoder.Decode(&data)
// 	if err != nil {
// 		LoglnArgs("Post File Error:", err)
// 		response := Response{Error: "error parsing request body"}
// 		w.WriteHeader(http.StatusBadRequest)
// 		WriteOutJSONMessage(response, w)
// 		return
// 	}
// 	defer req.Body.Close()
// 	headerObj := createHeaderObject(data.Attributes)
// 	filePath := bytes.NewBufferString(filepath.Join(CreateTodaysFolder(), string(data.ValidateFile)))
// 	n, err := sfile.WriteSaveFile(filePath.Bytes(), data.Data, headerObj, data.StartIndex, data.Size)
// 	if err != nil {
// 		log.Fatal(err)
// 		errReturn := map[string]interface{}{"Count": n, "Error": fmt.Sprintf("Error while writing file %s; %s", data.ValidateFile, err)}
// 		WriteOutJSONMessage(errReturn, w)
// 		return
// 	}
// 	Logf("File data, Name: %s. Wrote %d bytes", data.ValidateFile, n)
// 	WriteOutJSONMessage(map[string]interface{}{"Count": n, "Error": ""}, w)
// }

// ValidateFile is a GET request that takes in a file hash and checks to see
// if that file exists on the server as a whole.
func ValidateFile(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := Response{Error: "Only GET method allowed."}
		WriteOutJSONMessage(response, w)
		return
	}
	logger.Trace.ServerCall(req, "ValidateFile")
	errMsg := map[string]interface{}{"Error": "", "Size": 0}
	folder := req.URL.Query().Get("Folder")
	fileHash := req.URL.Query().Get("Hash")
	if fileHash == "" {
		index, err := strconv.Atoi(req.URL.Query().Get("Index"))
		if err != nil {
			errMsg["Error"] = err
			WriteOutJSONMessage(errMsg, w)
			return
		}
		validateFileWithIndex(w, req, folder, index)
	} else {
		validateFileWithHash(w, req, folder, fileHash)
	}
}

//
func validateFileWithIndex(w http.ResponseWriter, req *http.Request, folder string, index int) {

	Logf("validating %d from %s", index, folder)
	errMsg := map[string]interface{}{"Error": ""}
	files, err := ioutil.ReadDir(filepath.Join(rootPath, folder))
	if err != nil {
		errMsg["Error"] = err
		WriteOutJSONMessage(errMsg, w)
		return
	}
	if index >= len(files) || index < 0 {
		errMsg["Error"] = errors.New("error: index out of range")
		WriteOutJSONMessage(errMsg, w)
		return
	}
	saveFileObj, err := sfile.ReadSaveFile([]byte(files[index].Name()), nil)
	if err != nil {
		errMsg["Error"] = err
		WriteOutJSONMessage(errMsg, w)
		return
	}
	splitFilePath := strings.Split(files[index].Name(), string(os.PathSeparator))
	correctHash := []byte(splitFilePath[len(splitFilePath)-1])
	checkHash := sha256.Sum256(saveFileObj.Data)
	if bytes.Compare(correctHash, checkHash[:]) == 0 {
		WriteOutJSONMessage(errMsg, w)
		return
	}
	errMsg["Error"] = errors.New("error: original hash does not match current data hash")
	WriteOutJSONMessage(errMsg, w)
}

//
func validateFileWithHash(w http.ResponseWriter, req *http.Request, folder, hash string) {
	Logf("validating %s from %s", hash, folder)
	errMsg := map[string]interface{}{"Error": ""}
	saveFileObj, err := sfile.ReadSaveFile([]byte(filepath.Join(rootPath, folder, hash)), nil)
	if err != nil {
		errMsg["Error"] = err
		WriteOutJSONMessage(errMsg, w)
		return
	}
	correctHash := []byte(hash)
	checkHash := sha256.Sum256(saveFileObj.Data)
	if bytes.Compare(correctHash, checkHash[:]) == 0 {
		WriteOutJSONMessage(errMsg, w)
		return
	}
	errMsg["Error"] = errors.New("error: no file matches hash given")
	WriteOutJSONMessage(errMsg, w)
}

// GetFolders is a method to retrieve the list of folder names in the Data path.
func GetFolders(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := Response{Error: "Only GET method allowed."}
		WriteOutJSONMessage(response, w)
		return
	}
	LogServerCall(req, "GetFolders")
	// Grab all folders in Data directory
	foldersFromDir, err := ioutil.ReadDir(rootPath)
	if err != nil {
		LogFatal(err.Error())
		errFolders := Response{Error: "ERROR: Could not read Data directory."}
		WriteOutJSONMessage(errFolders, w)
		return
	}

	folders := FoldersList{Folders: make([]Folder, 0)}
	// Grab all folder info
	for _, obj := range foldersFromDir {
		// Grab files in folder
		subFiles, err := ioutil.ReadDir(filepath.Join(rootPath, obj.Name(), "."))
		if err != nil {
			LogFatal(err.Error())
			response := Response{Error: "error getting folders"}
			w.WriteHeader(http.StatusInternalServerError)
			WriteOutJSONMessage(response, w)
			return
		}
		// Create Folder object with its name and its file count
		f := Folder{Name: obj.Name(), Count: len(subFiles)}
		folders.Folders = append(folders.Folders, f)
	}
	response := Response{Data: folders, Error: ""}
	WriteOutJSONMessage(response, w)
}

// GetFiles is a method to accept a POST request for a specific folder in the Data path requesting files
// from startIndex to endIndex(exclusively). Must also send a dictionary with the keys of the attributes you want to extract
func GetFiles(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := Response{Error: "Only POST method allowed."}
		WriteOutJSONMessage(response, w)
		return
	}
	LogServerCall(req, "GetFiles")
	decoder := json.NewDecoder(req.Body)
	var data GetFilesWithAttributes
	err := decoder.Decode(&data)
	if err != nil {
		LoglnArgs("Post File Error:", err)
		response := Response{Error: "Could not decode the request body properly."}
		w.WriteHeader(http.StatusBadRequest)
		WriteOutJSONMessage(response, w)
		return
	}
	defer req.Body.Close()
	// check if out of range
	if data.StartIndex > data.EndIndex || data.StartIndex < 0 || data.EndIndex < 0 {
		LogFatal(fmt.Sprintf("ERROR: folder: %s, startIndex: %d, endIndex: %d", data.Folder, data.StartIndex, data.EndIndex))
		errFiles := Response{Error: fmt.Sprintf("ERROR: Either start or end index is incorrect. startIndex: %d, endIndex: %d", data.StartIndex, data.EndIndex)}
		WriteOutJSONMessage(errFiles, w)
		return
	}
	// get list of files from folder
	files, err := ioutil.ReadDir(filepath.Join(rootPath, data.Folder))
	if err != nil {
		log.Fatal(err)
		errFiles := Response{Error: "ERROR: Folder given could not be opened. Folder: " + data.Folder}
		WriteOutJSONMessage(errFiles, w)
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
			response := Response{Error: "Could not read file."}
			w.WriteHeader(http.StatusInternalServerError)
			WriteOutJSONMessage(response, w)
			return
		}
		// base64 encode data
		dstData := make([]byte, base64.StdEncoding.EncodedLen(len(saveFileObj.Data)))
		base64.StdEncoding.Encode(dstData, saveFileObj.Data)
		// map our objects
		headerMap := saveFileObj.Header.GetHeader()
		for i, v := range data.SortedAttributeKeys() {
			data.Attributes[v] = headerMap[i]
		}
		// Split on path separator to grab file hash
		validFile := strings.Split(obj.Name(), string(os.PathSeparator))
		// base64 encode file hash
		dstValid := make([]byte, base64.StdEncoding.EncodedLen(len(validFile[len(validFile)-1])))
		base64.StdEncoding.Encode(dstValid, []byte(validFile[len(validFile)-1]))
		// create our object
		f := FileData{Data: dstData, Size: int64(saveFileObj.Size), StartIndex: 0, ValidateFile: dstValid, Attributes: data.Attributes}
		allFiles.Files = append(allFiles.Files, f)
	}
	response := Response{Data: allFiles, Error: ""}
	WriteOutJSONMessage(response, w)
}
