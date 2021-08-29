package dismap

import (
	"testing"
	"time"
)

func Test_RequestSample(t *testing.T) {
	req, err := MakeDefaultRequest("https://www.baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	s, err := RequestSample(req, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", s)
}

func Test_MakeCustomRequest(t *testing.T) {
	req, err := MakeCustomRequest("https://www.baidu.com", "GET", "/robots.txt", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	s, err := RequestSample(req, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", s)
}
