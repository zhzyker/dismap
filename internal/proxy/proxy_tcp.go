package proxy

import (
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/pkg/logger"
	"golang.org/x/net/proxy"
	"net"
	"strconv"
	"time"
)

func ConnProxyTcp(host string, port int, timeout int) (net.Conn, error) {
	target := net.JoinHostPort(host, strconv.Itoa(port))
	scheme, address, proxyUri, err := parse.ProxyParse()
	if logger.DebugError(err) {
		return nil, err
	}

	var conn net.Conn
	if proxyUri == "" {
		conn, err = net.DialTimeout("tcp", target, time.Duration(timeout)*time.Second)
		if logger.DebugError(err) {
			return nil, err
		}
	}
	if scheme == "http" {
		conn, err = net.DialTimeout("tcp", target, time.Duration(timeout)*time.Second)
		if logger.DebugError(err) {
			return nil, err
		}
	}
	if scheme == "socks5" {
		tcp, err := proxy.SOCKS5("tcp", address, nil, &net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			KeepAlive: 10 * time.Second,
		})
		if logger.DebugError(err) {
			logger.Error("Cannot initialize socks5 proxy")
			return nil, err
		}
		conn, err = tcp.Dial("tcp", target)
		if logger.DebugError(err) {
			return nil, err
		}

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