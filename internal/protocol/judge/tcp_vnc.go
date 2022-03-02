package judge

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
)

func TcpVNC(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`^RFB \d`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "vnc"
		return true
	}
	return false
}