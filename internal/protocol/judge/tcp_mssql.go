package judge

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/proxy"
)

func TcpMssql(result *model.Result) bool {
	conn, err := proxy.ConnProxyTcp(result.Host, result.Port, flag.Timeout)
	if err != nil {
		return false
	}

	msg := "\x12\x01\x00\x34\x00\x00\x00\x00\x00\x00\x15\x00\x06\x01\x00\x1b\x00\x01\x02\x00\x1c\x00\x0c\x03\x00\x28\x00\x04\xff\x08\x00\x01\x55\x00\x00\x02\x4d\x53\x53\x51\x4c\x53\x65\x72\x76\x65\x72\x00\x00\x00\x31\x32"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return false
	}

	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	var buffer [256]byte
	if bytes.Equal(reply[:], buffer[:]) {
		return false
	} else if hex.EncodeToString(reply[0:4]) != "04010025" {
		return false
	} else {
		result.Protocol = "mssql"
	}

	v, bo := getVersion(reply)

	if bo {
		result.IdentifyBool = true
		result.Identify = append(result.Identify, model.HintFinger{
			Name:    "mssql",
			Version: v,
		})
	}
	result.Banner = parse.ByteToStringParse1(reply)
	result.BannerB = reply
	return true
}

func getVersion(reply []byte) (string, bool) {
	m, err := strconv.ParseUint(hex.EncodeToString(reply[29:30]), 16, 32)
	if err != nil {
		return "", false
	}
	s, err := strconv.ParseUint(hex.EncodeToString(reply[30:31]), 16, 32)
	if err != nil {
		return "", false
	}
	r, err := strconv.ParseUint(hex.EncodeToString(reply[31:33]), 16, 32)
	if err != nil {
		return "", false
	}
	v := fmt.Sprintf("%d.%d.%d", m, s, r)
	return v, true
}
