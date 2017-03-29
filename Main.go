package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type FileData struct {
	// change Data type to []byte when using base64
	// []byte looks for base64 encoding otherwise throws error
	Data    string
	Section int8
	Name    string
}

func PostFile(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var data FileData
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Post File Error:", err)
	}
	defer req.Body.Close()
	fmt.Printf("File data, Name: %s Section: %d Data: %s", data.Name, data.Section, string(data.Data))
	w.Write([]byte("File data recieved."))
}

func main() {
	http.HandleFunc("/post_file", PostFile)
	http.ListenAndServe(":8080", nil)
}
