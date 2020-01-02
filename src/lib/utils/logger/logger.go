package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"time"
)

// fileLogger logs to the secondary storage file
var fileLogger *log.Logger

// terminalandFileLogger logs to the secondary storage file and terminal
var terminalandFileLogger *log.Logger

const (
	logFilePrefix = "bench-route-"
	logDirectory  = "br-logs"
)

func init() {
	file, err := setupLogger()
	if err != nil {
		panic(err)
	}

	generalWriter := io.Writer(file)
	terminalWriter := io.MultiWriter(os.Stdout, file)

	fileLogger = log.New(generalWriter, "LOG:\t", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	terminalandFileLogger = log.New(terminalWriter, "LOG:\t", log.Ldate|log.Lmicroseconds|log.Lshortfile)

}

// setupLogger sets up files for logging
func setupLogger() (fp *os.File, err error) {
	currTime := time.Now()
	currFileName := fmt.Sprint(logFilePrefix, currTime.Format("2006-01-02#15:04:05"), ".log")
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	homePath := user.HomeDir
	logDirectoryPath := homePath + "/" + logDirectory
	err = os.MkdirAll(logDirectoryPath, 0755)
	if err != nil {
		return nil, err
	}
	logFilePath := logDirectoryPath + "/" + currFileName
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// File logs to the secondary storage.
// *code*:
// "p" -> Println()
// "f" -> Fatalln()
func File(msg string, code string) {
	switch code {
	case "p":
		fileLogger.Println(msg)

	case "f":
		fileLogger.Fatalln(msg)

	case "pa":
		fileLogger.Panicln(msg)
	}
}

// Terminal logs to the secondary storage and the terminal.
// *code*:
// "p" -> Println()
// "f" -> Fatalln()
func Terminal(msg string, code string) {
	switch code {
	case "p":
		terminalandFileLogger.Println(msg)

	case "f":
		terminalandFileLogger.Fatalln(msg)

	case "pa":
		terminalandFileLogger.Panicln(msg)
	}
}
