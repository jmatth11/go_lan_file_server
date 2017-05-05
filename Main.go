package main

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"io/ioutil"
	"strconv"
	"crypto/sha1"
)

type FileData struct {
	Data    []byte
	ValidateFile  []byte
	ValidateChunk []byte
	Size 		int64
	Type		string
	Name    string
}

type FileDataList struct {
	Files []FileData
	Error string
}

type FoldersList struct {
	Folders []Folder
	Error string
}

type Folder struct {
	Name string
	Count int
}

/**
 * Method to create the file in the current dates folder.
 * @param FileData
 */
func SpawnNewFile(fd FileData) error {
	path := CreateTodaysFolder()
	// create file
	file, err := os.Create(path + "\\" + fd.Name + fd.Type)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// write the data to the file
	_, err = file.Write(fd.Data)
	if err != nil {
		log.Fatal(err)
		return err
	}
	file.Close()
}

/**
 * Method to create a folder with the current date as its name.
 * @return string  The folder path
 */
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

/**
 * Method to handle post request for a file to be saved to the server.
 */
func PostFile(w http.ResponseWriter, req *http.Request) {
	// create json decoder
	decoder := json.NewDecoder(req.Body)
	// TODO This will go away and get replaced by post object's ranges
	data_range_str := req.Header.Get("Content-Range")
	var data FileData
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Post File Error:", err)
	}
	defer req.Body.Close()
	// create the temp file. TODO change up how to figure out size because size will be different
	tmp_file := CreateTempFile(data.Name + data.Type, data.Size + data_range_str.Len() + 1)
	WriteToTempFile(tmp_file, data, data_range)
	ValidateTempFile(tmp_file, data)
	//SpawnNewFile(data)
	log.Printf("File data, Name: %s", data.Name)
	// TODO create struct to json.Marshal and send back with an error or not.
	w.Write([]byte("File data recieved.\n"))
}

func CreateTempFile(file_name string, file_size int64) *File {
	// Create file
	file, err := os.Create(file_name)
	if err != nil {
		log.Fatal(err)
	}
	// Allocate file with certain size
	file.Truncate(file_size)
	return file
}

func WriteToTempFile(file *File, fd FileData, data_range_str string) {
	// TODO this will change after creating new file format
	data_range := data_range_str.Split("-")
	start_index := strconv.Atoi(data_range[0])
	end_index := strconv.Atoi(data_range[1])
	if (start_index == 0) {
		// TODO Header for new formated file will be different
		start_index, err := file.Write([]byte(data_range_str + "\n"))
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = file.WriteAt(fd.Data, start_index)
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateTempFile(file *File, fd FileData) {
	// TODO This file change when new file format implemented
	file_info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	var read_data []byte = make([]byte, file_info.Size())
	n, err := File.Read(read_data)
	if err != nil {
		log.Fatal(err)
	}
	//tmp_file_hash := sha1.Sum()
}

/**
 * Method to retrieve the list of folder names in the Data path.
 */
func GetFolders(w http.ResponseWriter, req *http.Request) {
	// Grab all folders in Data directory
	files, err := ioutil.ReadDir("Data\\.")
	if err != nil {
		log.Fatal(err)
		err_folders := FoldersList{Error: "ERROR: Could not read Data directory."}
		writeOutJsonError(err_folders, w)
		return
	}

	folders := FoldersList{Folders:make([]Folder, 0)}
  // Grab all folder info
	for _, obj := range files {
		// Grab files in folder
		sub_files, err := ioutil.ReadDir("Data\\" + obj.Name() + "\\.")
		if err != nil {
			log.Fatal(err)
		}
		// Create Folder object with its name and its file count
		f := Folder{Name:obj.Name(), Count:len(sub_files)}
		folders.Folders = append(folders.Folders, f)
	}
	b, err := json.Marshal(folders)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}

/**
 * Method to accept a GET request for a specific folder in the Data path requesting files
 * from start_index to end_index(exclusively).
 */
func GetFiles(w http.ResponseWriter, req *http.Request) {
	// grab folder name
	 folder := req.URL.Query().Get("folder")
	 // grab the start and end range of what files to grab
	 start_index, _ := strconv.Atoi(req.URL.Query().Get("start_index"))
	 end_index, _ := strconv.Atoi(req.URL.Query().Get("end_index"))
	 // check if out of range
	 if start_index > end_index || start_index < 0 || end_index < 0 {
		 log.Fatal(fmt.Sprintf("ERROR: folder: %s, start_index: %d, end_index: %d", folder, start_index, end_index))
		 err_files := FileDataList{Error:fmt.Sprintf("ERROR: Either start or end index is incorrect. start_index: %d, end_index: %d", start_index, end_index)}
		 writeOutJsonError(err_files, w)
		 return
	 }
	 // get list of files from folder
	 files, err := ioutil.ReadDir("Data\\" + folder)
	 if err != nil {
		 log.Fatal(err)
		 err_files := FileDataList{Error:"ERROR: Folder given could not be opened. Folder: " + folder}
		 writeOutJsonError(err_files, w)
		 return
	 }

	 all_files := FileDataList{Files:make([]FileData, 0)}
	 for _, obj := range files[start_index : end_index] {

		 // TODO Will need to change whenever new file format is implemented
		 src_dat, err := ioutil.ReadFile(obj.Name())
		 if err != nil {
			 log.Fatal(err)
		 }
		 dst_data := make([]byte, base64.StdEncoding.EncodedLen(len(src_dat)))
		 base64.StdEncoding.Encode(dst_data, src_dat)
		 f := FileData{Name:obj.Name(), Data:dst_data}
		 all_files.Files = append(all_files.Files, f)
	 }
	 b, err := json.Marshal(all_files)
	 if err != nil {
		 log.Fatal(err)
	 }
	 w.Write(b)
}

/**
 * Method to take an object json.Marshal it and write it out
 * to the console and the reposewriter.
 * @param obj interface{} A struct value
 * @param w http.ResponseWriter 
 */
func writeOutJsonError(obj interface{}, w http.ResponseWriter) {
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
