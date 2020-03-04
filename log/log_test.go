package log

import (
	"testing"
)

func TestDebug(t *testing.T) {
	Debug("debug")
	Info("debug")
	Warn("debug")
	Warnf("%v", 1)
	Error("debug")
}
