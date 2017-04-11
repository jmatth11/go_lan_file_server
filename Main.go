package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"io/ioutil"
)

type FileData struct {
	Data    []byte
	Name    string
}

// rename to ItemTypeList
type FoldersList struct {
	Folders []Folder
}

// rename to ItemType
type Folder struct {
	Name string
	Count int
}

func SpawnNewFile(fd FileData) {
	path := CreateTodaysFolder()
	file, err := os.Create(path + "\\" + fd.Name)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(fd.Data)
	if err != nil {
		panic(err)
	}
	file.Close()
}

func CreateTodaysFolder() string {
	year, month, day := time.Now().Date()
	name := fmt.Sprintf("Data\\%d-%d-%d", year, month, day)
	_, err := os.Stat(name)
	if err != nil {
		os.Mkdir(name, 0666)
	}
	return name
}


func PostFile(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data FileData
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Post File Error:", err)
	}
	defer req.Body.Close()
	SpawnNewFile(data)
	fmt.Printf("File data, Name: %s Data: %s", data.Name, string(data.Data))
	w.Write([]byte("File data recieved.\n"))
}

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

func GetFiles(w http.ResponseWriter, req *http.Request) {
	 folder := req.URL.Query().Get("folder")
	 files, err := ioutil.ReadDir("Data\\" + folder)
	 if err != nil {
		 log.Fatal(err)
	 }
	 all_files := FoldersList{Folders:make([]Folder, 0)}
	 for _, obj := range files {
		 f := Folder{Name:obj.Name(), Count:0}
		 all_files.Folders = append(all_files.Folders, f)
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
