package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"time"
)

// FileLogger Logs into the secondary storage file
var FileLogger *log.Logger

// TerminalandFileLogger Logs into the secondary storage file and terminal
var TerminalandFileLogger *log.Logger

const (
	logFilePrefix = "bench-route-"
	logDirectory  = "br-logs"
)

func init() {
	fmt.Println("Initializing logger")
	currTime := time.Now()
	currFileName := fmt.Sprint(logFilePrefix, currTime.Format("2006-01-02#15:04:05"), ".log")
	user, err := user.Current()
	if err != nil {
		fmt.Printf("Cannot access current user data\n")
		return
	}

	homePath := user.HomeDir
	logDirectoryPath := homePath + "/" + logDirectory
	err = os.MkdirAll(logDirectoryPath, 0755)
	if err != nil {
		fmt.Printf("error creating log directory : %s\n", logDirectoryPath)
		return
	}
	logFilePath := logDirectoryPath + "/" + currFileName
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0444)
	if err != nil {
		fmt.Printf("error opening log file : %s\n", logFilePath)
		return
	}

	generalWriter := io.Writer(file)
	terminalWriter := io.MultiWriter(os.Stdout, file)

	FileLogger = log.New(generalWriter, "LOG:\t", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	TerminalandFileLogger = log.New(terminalWriter, "LOG:\t", log.Ldate|log.Lmicroseconds|log.Lshortfile)

}
