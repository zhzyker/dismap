package judge

import (
	"encoding/hex"
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
	"strings"
)

func TcpTelnet(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`(Telnet>|^BeanShell)`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "telnet"
		return true
	} else if strings.Contains(hex.EncodeToString(buff[0:2]), "fffb") {
		result["protocol"] = "telnet"
		return true
	}
	return false
}
