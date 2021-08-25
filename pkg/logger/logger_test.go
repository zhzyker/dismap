package logger

import (
	"testing"
)

func TestPrint(t *testing.T) {
	Info("123")
	Warn("123")
	Error("1234")
	Success("1234")
	Failed("1234")

	Error(LightRed("12") + " 45")
}
