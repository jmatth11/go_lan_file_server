package server

// logging file to wrap log for later when I create a custom logger

import (
	"encoding/json"
	"log"
	"net/http"
)

// LogServerCall is used to log messages within the app.
func LogServerCall(req *http.Request, funcName string) {
	// right now just logging direct ip.
	// later might want to add req.Header.Get("X-Forwarded-For") to get possible tail of ips.
	log.Printf("%s|%s|%s %s", req.Method, funcName, "directly from:", req.RemoteAddr)
}

// WriteOutJSONMessage is a method to take an object json.Marshal it and write it out
// to the console and the reposewriter.
// @param obj interface{} A struct value
// @param w http.ResponseWriter
func WriteOutJSONMessage(obj interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("writeOutJSONMessage: %s", string(b))
	w.Write(b)
}

// Logln logs the text given
// @param test string
func Logln(text string) {
	log.Println(text)
}

// LoglnArgs logs text events and handles args
// @param text string
// @param args ...interface{}
func LoglnArgs(text string, args ...interface{}) {
	log.Print(text)
	for _, a := range args {
		log.Printf(" %s", a)
	}
	log.Print("\n")
}

// Logf logs the text given with the given args placed in the string.
// @param text string
// @param args ...interface{}
func Logf(text string, args ...interface{}) {
	log.Printf(text, args...)
}

// LogFatal logs text as fatal message
// @param text string
func LogFatal(text string) {
	log.Fatal(text)
}
