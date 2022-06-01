package protocol

import (
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

/*
func getTls(result map[string]interface{}, host string, port int, timeout int) bool {
	b, err := get.TlsProtocol(host, port, timeout)
	if err == nil {
		result["type"] = "tls"
		result["status"] = "open"
		result["banner.byte"] = b
		return true
	}
	return false
}

func getTcp(result map[string]interface{}, host string, port int, timeout int) bool {
	b, err := get.TcpProtocol(host, port, timeout)
	if err == nil {
		result["type"] = "tcp"
		result["status"] = "open"
		result["banner.byte"] = b
		return true
	}
	return false
}

func getUdp(result map[string]interface{}, host string, port int, timeout int) bool {
	b, err := get.UdpProtocol(host, port, timeout)
	if err == nil {
		result["type"] = "udp"
		result["status"] = "open"
		result["banner.byte"] = b
		return true
	}
	return false
}

func getInfo(result map[string]interface{}, host string, port int, timeout int, pt string) {
	var udpPort = []int{53, 111, 123, 137, 138, 139, 12345}
	switch pt {
	case "":
		if getTls(result, host, port, timeout) {
			return
		}
		if getTcp(result, host, port, timeout) {
			return
		}
		if isContainInt(udpPort, port) == true {
			if getUdp(result, host, port, timeout) {
				return
			}
		}
	case "tls":
		if getTls(result, host, port, timeout) {
			return
		}
	case "tcp":
		if getTcp(result, host, port, timeout) {
			return
		}
	case "udp":
		if getUdp(result, host, port, timeout) {
			return
		}
	}
	result["type"] = pt
	result["status"] = "close"
	result["banner.byte"] = make([]byte, 256)
	return
}

*/

func setResult(host string, port int, Args map[string]interface{}) map[string]interface{} {
	//timeout := Args["FlagTimeout"].(int)
	//scheme := Args["FlagMode"].(string)
	//pt := Args["FlagType"].(string)
	var banner []byte
	result := map[string]interface{}{
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
	//(result, host, port, timeout, pt)
	return result
}

/*
func Discover(host string, port int, Args map[string]interface{}) map[string]interface{} {
	result := setResult(host, port, Args)
	if result["status"] != "open" {
		return result
	}
	result["banner.string"] = parse.ByteToStringParse1(result["banner.byte"].([]byte))

	if result["type"] == "tls" {
		if JudgeTls(result, Args) {
			return result
		}
	}
	if result["type"] == "tcp" {
		if JudgeTcp(result, Args) {
			return result
		}
	}
	if result["type"] == "udp" {
		if JudgeUdp(result, Args) {
			return result
		}
	}
	return result
}
*/

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
