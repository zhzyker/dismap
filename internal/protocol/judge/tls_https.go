package judge

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func TlsHTTPS(result *model.Result) bool {
	var buff = result.BannerB
	ok, err := regexp.Match(`^HTTP/\d.\d \d*`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "https"
		httpResult, fpHints, httpErr := httpIdentifyResult(result)
		if logger.DebugError(httpErr) {
			result.Identify = fpHints
			result.Banner = "None"
			return true
		}
		result.Identify = fpHints
		result.Banner = httpResult.Title
		u, err := url.Parse(httpResult.Url)
		if err != nil {
			result.Path = ""
		} else {
			result.Path = u.Path
		}
		r := httpResult.Result
		c := fmt.Sprintf("[%s]", logger.Purple(httpResult.StatusCode))
		result.IdentifyBool = true
		result.Note = httpResult.Url
		if len(r) != 0 {
			result.IdentifyStr = fmt.Sprintf("%s %s", c, r)
			return true
		} else {
			result.IdentifyStr = c
			return true
		}
	}
	return false
}
