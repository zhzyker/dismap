package judge

import (
	"bytes"
	"encoding/hex"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/proxy"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpRMI(result *model.Result) bool {
	timeout := flag.Timeout
	host := result.Host
	port := result.Port

	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if logger.DebugError(err) {
		return false
	}

	msg := "\x4a\x52\x4d\x49\x00\x02\x4b"
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
	} else if hex.EncodeToString(reply[0:1]) != "4e" {
		return false
	}
	result.Protocol = "rmi"
	result.Banner = parse.ByteToStringParse1(reply)
	result.BannerB = reply
	return true
}
