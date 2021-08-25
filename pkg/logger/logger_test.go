package logger

import (
	"testing"
)

func all() {
	// Fatal("1234")
	Error("1234")
	Info("123")
	Success("1234")
	Failed("1234")

	Warn("123")
	Debug("123")
	Verbose("1234")
}
func TestPrint(t *testing.T) {
	t.Log("default level")
	all()

	t.Log("set level: verbose")
	SetLevel(LevelVerbose)
	all()

	t.Log("set level: error")
	SetLevel(LevelError)
	all()
}

func TestSetLevel(t *testing.T) {

}
