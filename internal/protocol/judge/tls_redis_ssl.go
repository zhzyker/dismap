package judge

import (
	"regexp"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TlsRedisSsl(result *model.Result) bool {
	var buff = result.BannerB
	ok, err := regexp.Match(`(^-ERR(.*)command|^-(.*).Redis)`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "redis-ssl"
		return true
	}
	return false
}
