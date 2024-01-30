package judge

import (
	"bytes"
	"encoding/hex"
	"strings"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/proxy"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpLDAP(result *model.Result) bool {
	timeout := flag.Timeout
	host := result.Host
	port := result.Port

	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if logger.DebugError(err) {
		return false
	}

	msg := "\x30\x0c\x02\x01\x01\x60\x07\x02\x01\x03\x04\x00\x80\x00"
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
	}
	if strings.Contains(hex.EncodeToString(reply), "010004000400") == false {
		return false
	}
	result.Protocol = "ldap"
	result.Banner = parse.ByteToStringParse2(reply[0:16])
	result.BannerB = reply
	return true
}
