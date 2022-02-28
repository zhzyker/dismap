package parse

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"net/url"
	"strconv"
)


func UriParse(target string) (string, string, string, int, error) {
	u, err := url.Parse(target)
	if logger.DebugError(err) {
		logger.Error(logger.LightRed(target) + " is not in uri format")
		return target, "", "", 0, err
	}
	s := u.Scheme
	p := u.Port()
	if s == "http" {
		if p == "" { p = "80" }
	}
	if s == "https" {
		if p == "" { p = "443" }
	}
	uPort, err := strconv.Atoi(p)
	if p == "" {
		logger.Error(logger.LightRed(target) + " is not in uri format, no port available")
		return target, "", "", 0, err
	}
	return target, s, u.Hostname(), uPort, nil
}
