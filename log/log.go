package log

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	debugPrefix = "[DEBUG]"
	infoPrefix  = "[INFO]"
	warnPrefix  = "[WARN]"
	errorPrefix = "[ERROR]"
)

func log(lvl, msg string) {
	_, f, l, ok := runtime.Caller(2)
	if !ok {
		fmt.Fprintln(os.Stderr, "something error occured.")
		return
	}

	now := time.Now().Format("2006/01/02 03:04:05")
	text := fmt.Sprintf("%s %s %s:%d %s", now, lvl, f, l, msg)
	fmt.Fprintln(os.Stderr, text)
}

func Debug(msg string) {
	log(debugPrefix, msg)
}

func Info(msg string) {
	log(infoPrefix, msg)
}

func Warn(msg string) {
	log(warnPrefix, msg)
}

func Warnf(format, msg string) {
	text := fmt.Sprintf(format, msg)
	log(warnPrefix, text)
}

func Error(msg string) {
	log(errorPrefix, msg)
}
