package judge

import (
	"regexp"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpIMAP(result *model.Result) bool {
	var buff []byte
	buff = result.BannerB
	ok, err := regexp.Match(`^* OK`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "imap"
		return true
	}
	return false
}
