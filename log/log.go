package log

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/jiro4989/textimg/v3/internal/global"
)

const (
	debugPrefix = "[DEBUG]"
	infoPrefix  = "[INFO]"
	warnPrefix  = "[WARN]"
	errorPrefix = "[ERROR]"
)

func log(lvl string, msg interface{}) {
	_, f, l, ok := runtime.Caller(2)
	if !ok {
		fmt.Fprintln(os.Stderr, "something error occured.")
		return
	}

	now := time.Now().Format("2006/01/02 03:04:05")
	text := fmt.Sprintf("%s %s %s %s:%d %v", now, global.AppName, lvl, f, l, msg)
	fmt.Fprintln(os.Stderr, text)
}

func Debug(msg interface{}) {
	log(debugPrefix, msg)
}

func Info(msg interface{}) {
	log(infoPrefix, msg)
}

func Warn(msg interface{}) {
	log(warnPrefix, msg)
}

func Warnf(format string, msg interface{}) {
	text := fmt.Sprintf(format, msg)
	log(warnPrefix, text)
}

func Error(msg interface{}) {
	log(errorPrefix, msg)
}
