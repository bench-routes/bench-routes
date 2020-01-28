package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
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

	fileLogger = log.New(generalWriter, "LOG:\t", log.Ldate|log.Lmicroseconds)
	terminalandFileLogger = log.New(terminalWriter, "LOG:\t", log.Ldate|log.Lmicroseconds)
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

func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}

func logLocater(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	_, file, line, _ := runtime.Caller(depth)
	location := chopPath(file) + ":" + strconv.Itoa(line)
	return location
}

// File logs to the secondary storage.
// *code*:
// "p" -> Println()
// "f" -> Fatalln()
// "pa" -> Panicln()
func File(msg string, code string) {
	switch code {
	case "p":
		logMessage := logLocater(2) + " " + msg
		fileLogger.Println(logMessage)

	case "f":
		logMessage := logLocater(2) + " " + msg
		fileLogger.Fatalln(logMessage)

	case "pa":
		logMessage := logLocater(2) + " " + msg
		fileLogger.Panicln(logMessage)
	}
}

// Terminal logs to the secondary storage and the terminal.
// *code*:
// "p" -> Println()
// "f" -> Fatalln()
// "pa" -> Panicln()
func Terminal(msg string, code string) {
	switch code {
	case "p":
		logMessage := logLocater(2) + " " + msg
		terminalandFileLogger.Println(logMessage)

	case "f":
		logMessage := logLocater(2) + " " + msg
		terminalandFileLogger.Fatalln(logMessage)

	case "pa":
		logMessage := logLocater(2) + " " + msg
		terminalandFileLogger.Panicln(logMessage)
	}
}
