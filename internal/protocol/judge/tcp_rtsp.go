package judge

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/proxy"
)

func TcpRTSP(result *model.Result) bool {
	ok, err := regexp.Match(`^RTSP/`, result.BannerB)
	if err != nil {
		return false
	}
	if ok {
		result.Protocol = "rtsp"
		return true
	}

	if rtsp(result) {
		return true
	}
	return false
}

func rtsp(result *model.Result) bool {
	address := net.JoinHostPort(result.Host, strconv.Itoa(result.Port))
	conn, err := proxy.ConnProxyTcp(result.Host, result.Port, flag.Timeout)
	if err != nil {
		return false
	}

	msg := fmt.Sprintf("OPTIONS rtsp://%s RTSP/1.0\r\nCSeq:1\r\n\r\n", address)
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
	} else if hex.EncodeToString(reply[0:4]) != "52545350" {
		return false
	}
	result.Protocol = "rtsp"
	result.Banner = parse.ByteToStringParse1(reply)
	result.BannerB = reply
	return true
}
