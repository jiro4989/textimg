package log

import "log"

const (
	debugPrefix = "[DEBUG]"
	infoPrefix  = "[INFO]"
	warnPrefix  = "[WARN]"
	errorPrefix = "[ERROR]"
)

func init() {
	log.SetPrefix("textimg: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Debug(msg string) {
	log.Println(debugPrefix, msg)
}

func Info(msg string) {
	log.Println(infoPrefix, msg)
}

func Warn(msg string) {
	log.Println(warnPrefix, msg)
}

func Warnf(format, msg string) {
	log.Printf(warnPrefix+format, msg)
}

func Error(msg string) {
	log.Println(errorPrefix, msg)
}
