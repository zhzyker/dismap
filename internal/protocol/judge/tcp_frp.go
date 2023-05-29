package judge

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/proxy"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpFrp(result *model.Result) bool {
	timeout := flag.Timeout
	host := result.Host
	port := result.Port

	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if logger.DebugError(err) {
		return false
	}

	msg := "\x00\x01\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00"
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
	} else if hex.EncodeToString(reply[0:12]) != "000100020000000100000000" {
		return false
	}
	result.Protocol = "frp"
	result.Banner = frpByteToStringParse(reply[0:12])
	result.BannerB = reply
	return true
}

func frpByteToStringParse(p []byte) string {
	var w []string
	var res string
	for i := 0; i < len(p); i++ {
		asciiTo16 := fmt.Sprintf("\\x%s", hex.EncodeToString(p[i:i+1]))
		w = append(w, asciiTo16)
	}
	res = strings.Join(w, "")
	return res
}
