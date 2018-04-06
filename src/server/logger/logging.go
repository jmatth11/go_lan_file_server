package logger

import (
	"io"
	"log"
	"net/http"
)

// ServLog is a wrapper type around log.Logger to provide extra logging functions
type ServLog struct {
	*log.Logger
}

var (
	// Trace level logging
	Trace *ServLog
	// Info level logging
	Info *ServLog
	// Warn level logging
	Warn *ServLog
	// Error level logging
	Error *ServLog
)

// InitLoggers Initializes the four available loggers
func InitLoggers(trace io.Writer, info io.Writer, warn io.Writer, err io.Writer) {
	Trace = &ServLog{log.New(trace, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)}
	Info = &ServLog{log.New(info, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)}
	Warn = &ServLog{log.New(warn, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)}
	Error = &ServLog{log.New(err, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)}
}

// ServerCall is a convience method for logging requests on the server
func (sl *ServLog) ServerCall(req *http.Request, funcName string) {
	sl.Printf("%s|%s|%s %s", req.Method, funcName, "directly from:", req.RemoteAddr)
}
