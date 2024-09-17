package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{}
}

func (LoggerObject *Logger) GetLoggerObject(infoFilePath, errorFilePath, layer string) *Logger {
	file, err := os.OpenFile(infoFilePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln("Error\t", err)
	}
	LoggerObject.InfoLogger = log.New(file, fmt.Sprint(layer)+" "+"INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	file, err = os.OpenFile(errorFilePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln("Error\t", err)
	}
	LoggerObject.ErrorLogger = log.New(file, fmt.Sprint(layer)+" "+"ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return LoggerObject
}
