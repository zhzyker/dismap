package judge

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
)

func TcpFTP(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`(^220(.*FTP|.*FileZilla)|^421(.*)connections)`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "ftp"
		return true
	}
	return false
}