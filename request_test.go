package dismap

import (
	"testing"
	"time"
)

func Test_RequestSample(t *testing.T) {
	s, err := RequestSample("https://www.baidu.com", 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", s)
}
