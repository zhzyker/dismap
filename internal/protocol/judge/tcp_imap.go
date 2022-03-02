package judge

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
)

func TcpIMAP(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`^* OK`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "imap"
		return true
	}
	return false
}