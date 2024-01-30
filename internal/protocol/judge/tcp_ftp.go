package judge

import (
	"regexp"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpFTP(result *model.Result) bool {
	var buff = result.BannerB
	ok, err := regexp.Match(`(^220(.*FTP|.*FileZilla)|^421(.*)connections)`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "ftp"
		return true
	}
	return false
}
