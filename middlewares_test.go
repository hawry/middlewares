package middlewares

import (
	"io"
	"os"
)

func ExampleSetOutput() {
	//This will print the log to logFile
	logFile, _ := os.Open("logfile")
	SetOutput(logFile)
}

func ExampleSetOutput_multiWriter() {
	//This will print the log both to Stdout and a file. Any ínterface implementing io.Writer can be used.
	logFile, _ := os.Open("logfile")
	mw := io.MultiWriter(os.Stdout, logFile)
	SetOutput(mw)
}
