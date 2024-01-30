package judge

import (
	"regexp"
	"strings"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpSSH(result *model.Result) bool {
	var buff []byte
	buff = result.BannerB
	ok, err := regexp.Match(`^SSH.\d`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		str := result.Banner
		result.Banner = strings.Split(str, "\\x0d\\x0a")[0]
		result.Protocol = "ssh"
		return true
	}
	return false
}
