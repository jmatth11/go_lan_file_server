package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type FileData struct {
	Data    []byte
	Name    string
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
	name := fmt.Sprintf("%d-%d-%d", year, month, day)
	_, err := os.Stat(name)
	if err != nil {
		os.Mkdir(name, 0666)
	}
	return name
}

// package mime/multipart. try using to help with large files
// example https://play.golang.org/p/MrE9BwNbB1
// stackoverflow source: http://stackoverflow.com/questions/20765859/go-accepting-http-post-multipart-files
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

func main() {
	http.HandleFunc("/post_file", PostFile)
	http.ListenAndServe(":8080", nil)
}
