package function

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var logFile *os.File

var GlobalId string

var turnOff = false //turn off log

func InitLogFileLin() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := "/logs/"
	currentTime := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("%slogfile_%s.txt", filePath, currentTime)

	// Create the directory if it doesn't exist (only for relative paths)
	if !filepath.IsAbs(filePath) {
		err := os.MkdirAll(filePath, 0755)
		if err != nil {
			return err
		}
	}

	// Print the resolved filepath for debugging
	filepath.Join(wd, fileName)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	logFile = file
	return nil
}

func InitLogFileWin() error {

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println("Working Directory:", wd)

	// Define the filepath
	filePath := "logs\\" // Using Windows filepath separator
	currentTime := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("%slogfile_%s.txt", filePath, currentTime)

	// Print the resolved filepath for debugging
	resolvedFilePath := filepath.Join(wd, fileName)

	// Create the directory if it doesn't exist
	err = os.MkdirAll(filepath.Join(wd, filePath), 0755)
	if err != nil {
		return err
	}

	// Open or create the file
	file, err := os.OpenFile(resolvedFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	logFile = file

	return nil
}

func GenerateRandomID() {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < 10; i++ {
		result += string(charset[rand.Intn(len(charset))])
	}

	ms := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	ms = ms[len(ms)-6:]

	result += time.Now().Format("20060102150405") + ms

	GlobalId = result
}

func GenerateAccountNumber(account *string) {
	ms := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	ms = ms[len(ms)-7:]

	result := time.Now().Format("20060102") + ms

	PrintLog(fmt.Sprintf("ACCOUNT NUMBER =  %s", result))

	*account = result
}

func PrintLog(detail string) {
	os := runtime.GOOS

	if turnOff {
		fmt.Println(detail)
		return
	}

	switch os {
	case "windows":
		if err := InitLogFileWin(); err != nil {
			fmt.Println("Error initializing log file:", err)
			return
		}
	case "linux":
		if err := InitLogFileLin(); err != nil {
			fmt.Println("Error initializing log file:", err)
			return
		}
	}

	if logFile != nil {
		defer logFile.Sync() // Make sure logs are written before program exit
		pc, file, line, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()

		// Extract only function name without package path
		lastSlashIndex := strings.LastIndex(funcName, "/")
		if lastSlashIndex >= 0 {
			funcName = funcName[lastSlashIndex+1:]
		}

		_, fileName := filepath.Split(file)
		logFile.WriteString(fmt.Sprintf("%s:%d:%s:%s: %s\n", strings.ToUpper(fileName), line, GlobalId, strings.ToUpper(funcName), detail))
	}

}
