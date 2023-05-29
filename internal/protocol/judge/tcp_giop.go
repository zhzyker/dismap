package judge

import (
	"encoding/hex"
	"strings"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/proxy"
)

func TcpGIOP(result *model.Result) bool {
	conn, err := proxy.ConnProxyTcp(result.Host, result.Port, flag.Timeout)
	if err != nil {
		return false
	}

	msg := "\x47\x49\x4f\x50\x01\x02\x00\x03\x00\x00\x00\x17\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00\x0b\x4e\x61\x6d\x65\x53\x65\x72\x76\x69\x63\x65"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return false
	}

	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	if !strings.Contains(hex.EncodeToString(reply[0:4]), "47494f50") {
		return false
	}

	result.Protocol = "giop"
	result.Banner = parse.ByteToStringParse2(reply[0:4])
	result.BannerB = reply
	return true
}
