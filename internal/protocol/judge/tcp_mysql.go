package judge

import (
	"regexp"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpMysql(result *model.Result) bool {
	var buff []byte
	buff = result.BannerB
	ok, err := regexp.Match(`(mysql_native_password|MySQL server|MariaDB server|mysqladmin flush-hosts)`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "mysql"
		return true
	}
	return false
}
