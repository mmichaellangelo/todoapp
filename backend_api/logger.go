package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

type LoggerMiddleware struct {
	handler http.Handler
	logfile *os.File
}

func NewLoggerMiddleware(handler http.Handler) *LoggerMiddleware {
	filename := time.Now().Format("2006-01-02_15-04-05") + ".txt"
	filepath := "logs/" + filename
	// create logs directory if not exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatalf("failed to create directories: %v", err)
	}
	// open logfile, create if not exists
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to open or create file: %v", err)
	}

	return &LoggerMiddleware{handler: handler, logfile: file}
}

func (h *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log the request
	h.logfile.WriteString(time.Now().Format(time.RFC3339) + " " + r.RemoteAddr + " " + r.Method + " " + r.URL.Path + "\n")

	h.handler.ServeHTTP(w, r)
}

func (h *LoggerMiddleware) Close() {
	h.logfile.Close()
}
