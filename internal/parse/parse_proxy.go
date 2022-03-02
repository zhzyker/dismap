package parse

import (
	"fmt"
	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/pkg/logger"
	"net/url"
)

func ProxyParse() (string, string, string, error) {
	if flag.Proxy == "" {
		return flag.Proxy, flag.Proxy, flag.Proxy, nil
	}

	u, err := url.Parse(flag.Proxy)
	if logger.DebugError(err) {
		logger.Fatal(fmt.Sprintf("The proxy address %s is formatted incorrectly", logger.LightRed(flag.Proxy)))
		return flag.Proxy, flag.Proxy, flag.Proxy, err
	}
	if u.Scheme == "http" || u.Scheme == "socks5" {
		return u.Scheme, u.Host, flag.Proxy, nil
	}
	logger.Fatal(logger.Red("Unsupported proxy protocol"))
	return flag.Proxy, flag.Proxy, flag.Proxy, nil
}
