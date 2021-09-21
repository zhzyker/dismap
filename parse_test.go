package dismap

import (
	"testing"
)

func TestParsePorts(t *testing.T) {
	t.Log(ParsePorts(""))
	t.Log(ParsePorts("1-23,80,8080,9000"))
}

func TestParseIPRange(t *testing.T) {
	t.Log(ParseIPRange("192.168.1.0/40"))
	t.Log(ParseIPRange("192.168.1.0/32"))
	t.Log(ParseIPRange("192.168.1.1-10"))
	t.Log(ParseIPRange("192.168.1-10.1-5"))
}
