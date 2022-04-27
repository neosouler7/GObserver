// Package alog provides async logging with goroutine.
package alog

import (
	"os"
	"time"
)

var (
	LogC        chan string
	logFileName = "gobserver.log"
)

// http://golang.site/go/article/210-%EC%B1%84%EB%84%90%EC%9D%84-%EC%9D%B4%EC%9A%A9%ED%95%9C-%EB%B9%84%EB%8F%99%EA%B8%B0-%EB%A1%9C%EA%B9%85
// TODO. golant basic log package + cycle.

// Initialize and returns logFile & channel.
func initALog() chan string {
	// Creates log file if not exists.
	if _, err := os.Stat(logFileName); os.IsNotExist(err) {
		f, _ := os.Create(logFileName)
		f.Close()
	}

	// Returns logging channel.
	return make(chan string)
}

// Listen for log.
func listenLog(logC chan string) {
	// Logs until channel is closed.
	for msg := range logC {
		f, _ := os.OpenFile(logFileName, os.O_APPEND, os.ModeAppend)
		f.WriteString(time.Now().String() + " " + msg + "\n")
		f.Close()
	}
}

// Starts alog package.
func Start() {
	LogC := initALog()
	go listenLog(LogC)
}
