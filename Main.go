package main

import (
	"net/http"
	"os"
	"server"
	"server/logger"
)

func main() {
	logger.InitLoggers(os.Stdout, os.Stdout, os.Stdout, os.Stdout)
	logger.Trace.Println("starting up server on port :8080")

	path := ""
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	handler := server.New(path)
	if err := http.ListenAndServe(":8080", handler); err != http.ErrServerClosed {
		logger.Error.Printf("%v\n", err)
	} else {
		logger.Trace.Println("server closed")
	}
}
