package main

import (
	"net/http"
	"log"
	"encoding/json"
	"strings"
	"os"
	"fmt"
)

type appError struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Error   error  `json:"-"`
	Message string `json:"message"`
}

/*
Encode error data to json and write it to http
 */
func (error *appError)Write(writer http.ResponseWriter) {
	if Cfg.LogToFile {
		dirPathSl := strings.Split(Cfg.LogFile, "/")
		dirPathSlLen := len(dirPathSl) - 1
		dirPath := strings.Join(dirPathSl[:dirPathSlLen], "/")
		
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			os.MkdirAll(dirPath, os.ModePerm)
		}
		
		logFile, err := os.OpenFile(Cfg.LogFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("error opening file: ", err)
		}
		defer logFile.Close()
		
		log.SetOutput(logFile)
		log.Println(error.Code, err)
	}
	writer.WriteHeader(error.Code)
	log.Printf("%s", error.Error)
	json.NewEncoder(writer).Encode(error)
}

func NewAppError(code int, err error) appError  {
	e := appError{
		Success: false,
		Code: code,
		Error: err,
		Message: fmt.Sprint(err),
	}
	return e
}