package judge

import (
	"regexp"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpRedis(result *model.Result) bool {
	var buff []byte
	buff = result.BannerB
	ok, err := regexp.Match(`(^-ERR(.*)command|^-DENIED.Redis)`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "redis"
		return true
	}
	return false
}
