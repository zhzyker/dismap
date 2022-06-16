package protocol

import (
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/protocol/get"
	"github.com/zhzyker/dismap/pkg/logger"
	"time"
)

func isContainInt(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func setResult(host string, port int, Args map[string]interface{}) map[string]interface{} {
	var banner []byte
	result := map[string]interface{}{
		"date":            time.Now().Unix(),
		"status":          "None",
		"banner.byte":     banner,
		"banner.string":   "None",
		"protocol":        Args["FlagMode"].(string),
		"type":            Args["FlagType"].(string),
		"host":            host,
		"port":            port,
		"uri":             "None",
		"note":            "None",
		"path":            "",
		"identify.bool":   false,
		"identify.string": "None",
	}
	return result
}

func DiscoverTls(host string, port int, Args map[string]interface{}) map[string]interface{} {
	result := setResult(host, port, Args)
	b, err := get.TlsProtocol(host, port, Args["FlagTimeout"].(int))
	if logger.DebugError(err) {
		return result
	}
	result["type"] = "tls"
	result["status"] = "open"
	result["banner.byte"] = b
	result["banner.string"] = parse.ByteToStringParse1(b)
	if JudgeTls(result, Args) {
		return result
	}
	return result
}

func DiscoverTcp(host string, port int, Args map[string]interface{}) map[string]interface{} {
	result := setResult(host, port, Args)
	b, err := get.TcpProtocol(host, port, Args["FlagTimeout"].(int))
	if logger.DebugError(err) {
		return result
	}
	result["type"] = "tcp"
	result["status"] = "open"
	result["banner.byte"] = b
	result["banner.string"] = parse.ByteToStringParse1(b)
	if JudgeTcp(result, Args) {
		return result
	}
	return result
}

func DiscoverUdp(host string, port int, Args map[string]interface{}) map[string]interface{} {
	result := setResult(host, port, Args)
	var udpPort = []int{53, 111, 123, 137, 138, 139, 12345}
	if isContainInt(udpPort, port) {
		return result
	}
	b, err := get.UdpProtocol(host, port, Args["FlagTimeout"].(int))
	if logger.DebugError(err) {
		return result
	}
	result["type"] = "tcp"
	result["status"] = "open"
	result["banner.byte"] = b
	result["banner.string"] = parse.ByteToStringParse1(b)
	if JudgeUdp(result, Args) {
		return result
	}
	return result
}
