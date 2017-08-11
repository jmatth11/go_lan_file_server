package main

import (
	"net/http"
	"os"
	"server"
)

// PostFile is a method to handle post request for a file to be saved to the server.
func PostFile(w http.ResponseWriter, req *http.Request) {
	server.WriteFile(w, req)
}

// GetFolders is a method to retrieve the list of folder names in the Data path.
func GetFolders(w http.ResponseWriter, req *http.Request) {
	server.GetFolders(w, req)
}

// GetFiles is a method to accept a POST request for a specific folder in the Data path requesting files
// from startIndex to endIndex(exclusively). Must also send a dictionary with the keys of the attributes you want to extract
func GetFiles(w http.ResponseWriter, req *http.Request) {
	server.GetFiles(w, req)
}

// ValidateFile is a GET request that takes in a file hash and checks to see
// if that file exists on the server as a whole.
func ValidateFile(w http.ResponseWriter, req *http.Request) {
	server.ValidateFile(w, req)
}

// PingServ method listens for any message and sends back a response that lets
// the user know it is hitting the right address.
func PingServ(w http.ResponseWriter, req *http.Request) {
	server.PingServ(w, req)
}

func main() {
	// Check if Data Folder exists and if not, create it.
	server.Logln("starting up server on port :8080")
	_, err := os.Stat("Data")
	if err != nil {
		server.Logln("creating Initial Data folder")
		os.Mkdir("Data", 0666)
	}
	// register functions
	http.HandleFunc("/ping", PingServ)
	http.HandleFunc("/post_file", PostFile)
	http.HandleFunc("/get_folders", GetFolders)
	http.HandleFunc("/get_files", GetFiles)
	http.HandleFunc("/validate_file", ValidateFile)
	http.ListenAndServe(":8080", nil)
}
