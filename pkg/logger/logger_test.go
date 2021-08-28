package logger

import (
	"testing"
)

func all() {
	// Fatalf("name: %s\n", "err")
	// Fatalln("123", "456")

	Errorf("name: %s\n", "err")
	Errorln("123", "456")

	Infof("name: %s\n", "err")
	Infoln("123", "456")

	Successf("name: %s\n", "err")
	Successln("123", "456")

	Failedf("name: %s\n", "err")
	Failedln("123", "456")

	Warnf("name: %s\n", "err")
	Warnln("123", "456")

	Debugf("name: %s\n", "err")
	Debugln("123", "456")

	Verbosef("name: %s\n", "err")
	Verboseln("123", "456")
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
