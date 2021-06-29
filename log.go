package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	// UNSPECIFIED logs nothing
	UNSPECIFIED Level = iota // 0 :
	// TRACE logs everything
	TRACE // 1
	// INFO logs Info, Warnings and Errors
	INFO // 2
	// WARNING logs Warning and Errors
	WARNING // 3
	// ERROR just logs Errors
	ERROR // 4
)

// Level holds the log level.
type Level int

// Package level variables, which are pointers to log.Logger.
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// initLog initializes log.Logger objects
func initLog(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	isFlag bool) {

	// Flags for defining the logging properties, to log.New
	flag := 0
	if isFlag {
		flag = log.Ldate | log.Ltime | log.Lshortfile
	}

	// Create log.Logger objects.
	Trace = log.New(traceHandle, "TRACE: ", flag)
	Info = log.New(infoHandle, "INFO : ", flag)
	Warning = log.New(warningHandle, "WARN : ", flag)
	Error = log.New(errorHandle, "ERROR: ", flag)
}

// SetLogLevel sets the logging level preference
func SetLogLevel() (err error) {

	// Creates os.*File, which has implemented io.Writer interface
	logfile, err := os.OpenFile(config.SystemLog.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		Error.Println("Problems opening log file. ", err)
		return
	}

	// Calls function initLog by specifying log level preference.
	switch Level(config.SystemLog.Level) {
	case TRACE:
		initLog(logfile, logfile, logfile, logfile, true)
		return
	case INFO:
		initLog(ioutil.Discard, logfile, logfile, logfile, true)
		return
	case WARNING:
		initLog(ioutil.Discard, ioutil.Discard, logfile, logfile, true)
		return
	case ERROR:
		initLog(ioutil.Discard, ioutil.Discard, ioutil.Discard, logfile, true)
		return
	default:
		initLog(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard, false)
		defer logfile.Close()
		return
	}
}
