package judge

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
)

func TcpPOP3(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`^\+OK`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "pop3"
		return true
	}
	return false
}