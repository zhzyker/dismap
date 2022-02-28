package judge

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
)

func TcpMysql(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`(mysql_native_password|MySQL server|MariaDB server|mysqladmin flush-hosts)`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "mysql"
		return true
	}
	return false
}
