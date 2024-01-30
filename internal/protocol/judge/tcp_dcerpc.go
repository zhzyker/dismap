package judge

import (
	"bytes"
	"encoding/hex"
	"strings"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/proxy"
)

func TcpDceRpc(result *model.Result) bool {
	conn, err := proxy.ConnProxyTcp(result.Host, result.Port, flag.Timeout)
	if err != nil {
		return false
	}

	msg1 := "\x05\x00\x0b\x03\x10\x00\x00\x00\x48\x00\x00\x00\x01\x00\x00\x00\xf8\x0f\xf8\x0f\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x01\x00\xc4\xfe\xfc\x99\x60\x52\x1b\x10\xbb\xcb\x00\xaa\x00\x21\x34\x7a\x00\x00\x00\x00\x04\x5d\x88\x8a\xeb\x1c\xc9\x11\x9f\xe8\x08\x00\x2b\x10\x48\x60\x02\x00\x00\x00"
	msg2 := "\x05\x00\x00\x03\x10\x00\x00\x00\x18\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05\x00"
	_, err = conn.Write([]byte(msg1))
	if err != nil {
		return false
	}
	reply1 := make([]byte, 256)
	_, _ = conn.Read(reply1)

	if hex.EncodeToString(reply1[0:8]) != "05000c0310000000" {
		return false
	}

	_, err = conn.Write([]byte(msg2))
	if err != nil {
		return false
	}

	reply2 := make([]byte, 512)
	_, _ = conn.Read(reply2)
	if conn != nil {
		_ = conn.Close()
	}

	result.Protocol = "dcerpc"

	c := 0
	zero := make([]byte, 1)
	var buffer bytes.Buffer
	for i := 0; i < len(reply2[42:]); {
		b := reply2[42:][i : i+2]
		i += 2
		if 42+i == len(reply2[42:]) {
			break
		}
		if string(b) == "\x09\x00" {
			break
		}
		if string(b) == "\x07\x00" {
			c += 1
			if c == 6 {
				break
			}
			buffer.Write([]byte("\x7C\x7C"))
			result.Banner = strings.Join([]string{buffer.String()}, ",")
			continue
		}
		if bytes.Equal(b[0:1], zero[0:1]) {
			continue
		}
		buffer.Write(b[0:1])
		result.Banner = strings.Join([]string{buffer.String()}, ",")
		if c == 6 {
			break
		}
	}
	result.BannerB = reply2
	return true
}
