package judge

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/proxy"
	"github.com/zhzyker/dismap/pkg/logger"
	"regexp"
	"strconv"
)

func TcpSocks(result map[string]interface{}, Args map[string]interface{}) bool {
	timeout := Args["FlagTimeout"].(int)

	if socks4(result, timeout) {
		result["protocol"] = "socks4"
		return true
	}

	if socks5(result, timeout) {
		result["protocol"] = "socks5"
		return true
	}
	return false
}

func socks5(result map[string]interface{}, timeout int) bool {
	host := result["host"].(string)
	port := result["port"].(int)
	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if logger.DebugError(err) {
		return false
	}

	msgSocks5 := "\x05\x02\x00\x02"
	/*
		\x05 - Version: 5
		\x02 - Authentication Method Count: 2
		\x00 - Method[0]: 0 (No authentication)
		\x02 - Method[1]: 2 (Username/Password)
	*/
	_, err = conn.Write([]byte(msgSocks5))
	if logger.DebugError(err) {
		return false
	}

	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	var buffer bytes.Buffer
	if string(reply[0]) == "\x05" {
		buffer.WriteString(fmt.Sprintf("[%s]", logger.LightYellow("Version:Socks5")))
	} else {
		return false
	}
	if string(reply[1]) == "\x00" {
		buffer.WriteString(fmt.Sprintf("[%s]", logger.LightYellow("Method:No Authentication(\\x00)")))
	}
	if string(reply[1]) == "\x02" {
		buffer.WriteString(fmt.Sprintf("[%s]", logger.LightYellow("Method:Username/Password(\\x02)")))
	}
	result["identify.bool"] = true
	result["identify.string"] = buffer.String()
	result["banner.string"] = parse.ByteToStringParse2(reply[0:2])
	result["banner.byte"] = reply
	return true
}

func socks4(result map[string]interface{}, timeout int) bool {
	host := result["host"].(string)
	port := result["port"].(int)
	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if logger.DebugError(err) {
		return false
	}

	p1 := strconv.FormatInt(int64(port/256), 16)
	if p1str, _ := strconv.Atoi(p1); p1str < 10 {
		p1 = fmt.Sprintf("0%s",p1)
	}
	p1byte, err := hex.DecodeString(p1)
	if logger.DebugError(err) {
		return false
	}
	p2 := strconv.FormatInt(int64(port%256), 16)
	if p2str, _ := strconv.Atoi(p2); p2str < 10 {
		p2 = fmt.Sprintf("0%s",p2)
	}
	p2byte, err := hex.DecodeString(p2)
	if logger.DebugError(err) {
		return false
	}
	msgByte := []byte {0x04, 0x01}
	msgByte = append(msgByte, p1byte[0])
	msgByte = append(msgByte, p2byte[0])
	msgStr := hex.EncodeToString(msgByte)

	grep := regexp.MustCompile(`(\d*).(\d*).(\d*).(\d*)`)
	ip := grep.FindStringSubmatch(host)[1:5]
	for _, i := range ip {
		if i == "0" {
			msgStr += "00"
			continue
		}
		i64, _ := strconv.ParseInt(i, 10, 64)

		n := strconv.FormatInt(i64, 16)
		if len(n) != 2 {
			n = fmt.Sprintf("0%s",n)
		}
		msgStr += n
	}
	msgStr += "0100"
	hexData, _ := hex.DecodeString(msgStr)
	_, err = conn.Write(hexData)
	if logger.DebugError(err) {
		return false
	}

	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	if string(reply[1]) == "\x5b" {
		result["banner.string"] = parse.ByteToStringParse2(reply[0:8])
		result["banner.byte"] = reply
		return true
	}
	return false
}