package judge

import (
	"regexp"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpSMTP(result *model.Result) bool {
	var buff []byte
	buff = result.BannerB
	ok, err := regexp.Match(`(^220[ -](.*)ESMTP|^421(.*)Service not available|^554 )`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "smtp"
		return true
	}
	return false
}
