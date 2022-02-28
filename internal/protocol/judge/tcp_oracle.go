package judge

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/proxy"
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
	"strconv"
)

func TcpOracle(result map[string]interface{}, Args map[string]interface{}) bool {
	timeout := Args["FlagTimeout"].(int)
	host := result["host"].(string)
	port := result["port"].(int)

	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if logger.DebugError(err) {
		return false
	}

	msg := "\x00\x5a\x00\x00\x01\x00\x00\x00\x01\x36\x01\x2c\x00\x00\x08\x00\x7f\xff\x7f\x08\x00\x00\x00\x01\x00\x20\x00\x3a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x34\xe6\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x28\x43\x4f\x4e\x4e\x45\x43\x54\x5f\x44\x41\x54\x41\x3d\x28\x43\x4f\x4d\x4d\x41\x4e\x44\x3d\x56\x45\x52\x53\x49\x4f\x4e\x29\x29"
	_, err = conn.Write([]byte(msg))
	if logger.DebugError(err) {
		return false
	}

	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	ok, err := regexp.Match(`\(DESCRIPTION=`, result["banner.byte"].([]byte))
	if ok {
		result["protocol"] = "oracle"
	} else {
		var buffer [256]byte
		if bytes.Equal(reply[:], buffer[:]) {
			return false
		} else if hex.EncodeToString(reply[0:8]) != "0065000004000000" {
			return false
		} else {
			result["protocol"] = "oracle"
		}
	}
	var vsnnum string
	banStr := parse.ByteToStringParse2(reply)
	grep := regexp.MustCompile(`\(VSNNUM=(\d*)\)`)
	vsnnum = grep.FindStringSubmatch(banStr)[1]
	v, err := strconv.ParseInt(vsnnum, 10, 64)
	if logger.DebugError(err) {
		result["identify.bool"] = false
	}

	hexVsnnum := strconv.FormatInt(v, 16)
	maj, err := strconv.ParseUint(hexVsnnum[0:1], 16, 32)
	min, err := strconv.ParseUint(hexVsnnum[1:2], 16, 32)
	a, err := strconv.ParseUint(hexVsnnum[2:4], 16, 32)
	b, err := strconv.ParseUint(hexVsnnum[4:5], 16, 32)
	c, err := strconv.ParseUint(hexVsnnum[5:7], 16, 32)

	var version string
	if err == nil {
		version = fmt.Sprintf("%s.%s.%s.%s.%s",
			strconv.FormatUint(maj,10),
			strconv.FormatUint(min,10),
			strconv.FormatUint(a,10),
			strconv.FormatUint(b,10),
			strconv.FormatUint(c,10),
		)
	} else {
		result["identify.bool"] = false
	}
	result["identify.bool"] = true
	result["identify.string"] = fmt.Sprintf("[%s]", logger.LightYellow(fmt.Sprintf("Version:%s", version)))
	result["banner.string"] = banStr
	result["banner.byte"] = reply
	return true
}
