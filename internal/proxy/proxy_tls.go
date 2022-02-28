package proxy

import (
	"crypto/tls"
	"github.com/zhzyker/dismap/pkg/logger"
	"net"
	"strconv"
	"time"
)

func ConnProxyTls(host string, port int, timeout int) (net.Conn, error) {
	target := net.JoinHostPort(host, strconv.Itoa(port))
	// scheme, address, proxyUri, err := parse.ProxyParse()
	// TLS does not support proxy function temporarily
	// 2022-02-23 by zhzyker
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: time.Duration(timeout)*time.Second},
		"tcp",
		target,
		&tls.Config{InsecureSkipVerify: true})
	if logger.DebugError(err) {
		return nil, err
	}
	err = conn.SetDeadline(time.Now().Add(time.Duration(timeout)*time.Second))
	if logger.DebugError(err) {
		if conn != nil {
			_ = conn.Close()
		}
		return nil, err
	}
	return conn, nil
}