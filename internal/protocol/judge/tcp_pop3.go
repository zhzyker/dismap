package judge

import (
	"regexp"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpPOP3(result *model.Result) bool {
	var buff []byte
	buff = result.BannerB
	ok, err := regexp.Match(`^\+OK`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "pop3"
		return true
	}
	return false
}
