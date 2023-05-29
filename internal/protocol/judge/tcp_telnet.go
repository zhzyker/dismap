package judge

import (
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TcpTelnet(result *model.Result) bool {
	var buff = result.BannerB
	ok, err := regexp.Match(`(Telnet>|^BeanShell)`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "telnet"
		return true
	} else if strings.Contains(hex.EncodeToString(buff[0:2]), "fffb") {
		result.Protocol = "telnet"
		return true
	}
	return false
}
