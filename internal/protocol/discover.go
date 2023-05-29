package protocol

import (
	"time"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/protocol/get"
	"github.com/zhzyker/dismap/pkg/logger"
)

func isContainInt(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func NewResult(host string, port int) *model.Result {
	return &model.Result{
		Date: time.Now(),
		Host: host,
		Port: port,
	}
}

func DiscoverTls(host string, port int) *model.Result {
	result := NewResult(host, port)

	b, err := get.TlsProtocol(host, port, flag.Timeout)
	if logger.DebugError(err) {
		return result
	}
	result.Type = "tls"
	result.Status = "open"
	result.BannerB = b
	result.Banner = parse.ByteToStringParse1(b)
	if JudgeTls(result) {
		return result
	}
	return result
}

func DiscoverTcp(host string, port int) *model.Result {
	result := NewResult(host, port)
	b, err := get.TcpProtocol(host, port, flag.Timeout)
	if logger.DebugError(err) {
		return result
	}
	result.Type = "tcp"
	result.Status = "open"
	result.BannerB = b
	result.Banner = parse.ByteToStringParse1(b)
	if JudgeTcp(result) {
		return result
	}
	return result
}

func DiscoverUdp(host string, port int) *model.Result {
	result := NewResult(host, port)
	var udpPort = []int{53, 111, 123, 137, 138, 139, 12345}
	if isContainInt(udpPort, port) {
		return result
	}
	b, err := get.UdpProtocol(host, port, flag.Timeout)
	if logger.DebugError(err) {
		return result
	}
	result.Type = "udp"
	result.Status = "open"
	result.BannerB = b
	result.Banner = parse.ByteToStringParse1(b)
	if JudgeUdp(result) {
		return result
	}
	return result
}
