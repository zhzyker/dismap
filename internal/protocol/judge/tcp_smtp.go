package judge

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
)

func TcpSMTP(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`(^220[ -](.*)ESMTP|^421(.*)Service not available|^554 )`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "smtp"
		return true
	}
	return false
}