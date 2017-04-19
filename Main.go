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
	// ValidateFile  []byte
	// ValidateChunk []byte
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
	file, err := os.Create(path + "\\" + fd.Name + fd.Type)
	if err != nil {
		log.Fatal(err)
		return err
	}
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
	year, month, day := time.Now().Date()
	name := fmt.Sprintf("Data\\%d-%d-%d", year, month, day)
	_, err := os.Stat(name)
	if err != nil {
		os.Mkdir(name, 0666)
	}
	return name
}

/**
 * Method to handle post request for a file to be saved to the server.
 */
func PostFile(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	data_range_str := req.Header.Get("Content-Range")
	var data FileData
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Post File Error:", err)
	}
	defer req.Body.Close()
	tmp_file := CreateTempFile(data.Name + data.Type, data.Size + data_range_str.Len() + 1)
	WriteToTempFile(tmp_file, data, data_range)
	ValidateTempFile(tmp_file, data)
	//SpawnNewFile(data)
	log.Printf("File data, Name: %s", data.Name)
	w.Write([]byte("File data recieved.\n"))
}

func CreateTempFile(file_name string, file_size int64) *File {
	file, err := os.Create(file_name)
	if err != nil {
		log.Fatal(err)
	}
	file.Truncate(file_size)
	return file
}

func WriteToTempFile(file *File, fd FileData, data_range_str string) {
	data_range := data_range_str.Split("-")
	start_index := strconv.Atoi(data_range[0])
	end_index := strconv.Atoi(data_range[1])
	if (start_index == 0) {
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
	files, err := ioutil.ReadDir("Data\\.")
	if err != nil {
		log.Fatal(err)
	}
	folders := FoldersList{Folders:make([]Folder, 0)}
	for _, obj := range files {
		sub_files, err := ioutil.ReadDir("Data\\" + obj.Name() + "\\.")
		if err != nil {
			log.Fatal(err)
		}
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
	 folder := req.URL.Query().Get("folder")
	 start_index, _ := strconv.Atoi(req.URL.Query().Get("start_index"))
	 end_index, _ := strconv.Atoi(req.URL.Query().Get("end_index"))
	 if start_index > end_index {
		 log.Fatal(fmt.Sprintf("ERROR: folder: %s, start_index: %d, end_index: %d", folder, start_index, end_index))
		 err_files := FileDataList{Error:"ERROR: Start index is greater than End index"}
		 b, err := json.Marshal(err_files)
		 if err != nil {
			 log.Fatal(err)
		 }
		 w.Write(b)
		 return
	 }
	 files, err := ioutil.ReadDir("Data\\" + folder)
	 if err != nil {
		 log.Fatal(err)
	 }
	 all_files := FileDataList{Files:make([]FileData, 0)}
	 for _, obj := range files[start_index : end_index] {
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

func main() {
	_, err := os.Stat("Data")
	if err != nil {
		os.Mkdir("Data", 0666)
	}
	http.HandleFunc("/post_file", PostFile)
	http.HandleFunc("/get_folders", GetFolders)
	http.HandleFunc("/get_files", GetFiles)
	http.ListenAndServe(":8080", nil)
}
