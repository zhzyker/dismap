package dismap

import (
	"testing"
)

func TestPingHost(t *testing.T) {
	t.Log(pingHost("10.10.30.177", 2))
	t.Log(pingHost("10.10.30.110", 2))
}

func Test_isDomainName(t *testing.T) {
	t.Log(isDomainName("www.baidu.com."))
	t.Log(isDomainName("12www.baidu.com/"))
	t.Log(isDomainName("www.baidu.com"))
}

func Test_CreateFile(t *testing.T) {
	opts := &Options{
		OutPut: "output.test.txt",
	}
	r := NewRunner(opts)
	file, err := r.parseOutputFile()
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	file.WriteString("testing \n")
}
