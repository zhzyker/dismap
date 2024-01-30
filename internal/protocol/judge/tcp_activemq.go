package judge

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"

	"github.com/zhzyker/dismap/internal/model"
)

func TcpActiveMQ(result *model.Result) bool {
	ok, err := regexp.Match(`ActiveMQ`, result.BannerB)
	if err != nil {
		return false
	}
	if ok {
		ver, err := strconv.ParseUint(hex.EncodeToString(result.BannerB[13:17]), 16, 32)
		if err == nil {
			version := fmt.Sprintf("Version:%s", strconv.FormatUint(ver, 10))
			result.IdentifyBool = true
			result.Identify = append(result.Identify, model.HintFinger{Name: "ActiveMQ", Version: version})
		}
		result.Protocol = "activemq"
		return true
	}
	return false
}
