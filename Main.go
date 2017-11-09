package main

import (
	"net/http"
	"os"
	"server"
)

func main() {
	// Check if Data Folder exists and if not, create it.
	server.Logln("starting up server on port :8080")
	path := ""
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	handler := server.New(path)
	if err := http.ListenAndServe(":8080", handler); err != http.ErrServerClosed {
		server.Logf("%v", err)
	} else {
		server.Logln("server closed")
	}
}
