package judge

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
	"strings"
)

func TcpSSH(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`^SSH.\d`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		str := result["banner.string"].(string)
		result["banner.string"] = strings.Split(str, "\\x0d\\x0a")[0]
		result["protocol"] = "ssh"
		return true
	}
	return false
}