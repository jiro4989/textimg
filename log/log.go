package log

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/jiro4989/textimg/v3/config/version"
)

const (
	debugPrefix = "[DEBUG]"
	infoPrefix  = "[INFO]"
	warnPrefix  = "[WARN]"
	errorPrefix = "[ERROR]"
)

func log(lvl string, msg any) {
	_, f, l, ok := runtime.Caller(2)
	if !ok {
		fmt.Fprintln(os.Stderr, "something error occurred.")
		return
	}

	now := time.Now().Format("2006/01/02 03:04:05")
	text := fmt.Sprintf("%s %s %s %s:%d %v", now, version.AppName, lvl, f, l, msg)
	fmt.Fprintln(os.Stderr, text)
}

func Debug(msg any) {
	log(debugPrefix, msg)
}

func Info(msg any) {
	log(infoPrefix, msg)
}

func Warn(msg any) {
	log(warnPrefix, msg)
}

func Warnf(format string, msg any) {
	text := fmt.Sprintf(format, msg)
	log(warnPrefix, text)
}

func Error(msg any) {
	log(errorPrefix, msg)
}
