package judge

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/proxy"
	"github.com/zhzyker/dismap/pkg/logger"
	"strconv"
)

func TcpMssql(result map[string]interface{}, Args map[string]interface{}) bool {
	timeout := Args["FlagTimeout"].(int)
	host := result["host"].(string)
	port := result["port"].(int)

	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if logger.DebugError(err) {
		return false
	}

	msg := "\x12\x01\x00\x34\x00\x00\x00\x00\x00\x00\x15\x00\x06\x01\x00\x1b\x00\x01\x02\x00\x1c\x00\x0c\x03\x00\x28\x00\x04\xff\x08\x00\x01\x55\x00\x00\x02\x4d\x53\x53\x51\x4c\x53\x65\x72\x76\x65\x72\x00\x00\x00\x31\x32"
	_, err = conn.Write([]byte(msg))
	if logger.DebugError(err) {
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
		result["protocol"] = "mssql"
	}

	v, bo := getVersion(reply)

	if bo {
		result["identify.bool"] = true
		result["identify.string"] = fmt.Sprintf("[%s]", logger.LightYellow(fmt.Sprintf("Version:%s", v)))

	}
	result["banner.string"] = parse.ByteToStringParse1(reply)
	result["banner.byte"] = reply
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