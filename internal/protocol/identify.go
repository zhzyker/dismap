package protocol

import (
	"bytes"
	"fmt"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/protocol/judge"
	"github.com/zhzyker/dismap/pkg/logger"
)

func JudgeTcp(result map[string]interface{}, Args map[string]interface{}) bool {
	protocol := result["protocol"].(string)
	runAll := true
	if protocol != "" {
		runAll = false
	}
	if protocol == "http" || protocol == "https" || runAll {
		if judge.TcpHTTP(result, Args) {
			printSuccess("TCP/HTTP", result)
			return true
		}
	}
	if protocol == "mysql" || runAll {
		if judge.TcpMysql(result) {
			printSuccess("TCP/Mysql", result)
			return true
		}
	}
	if protocol == "redis" || runAll {
		if judge.TcpRedis(result) {
			printSuccess("TCP/Redis", result)
			return true
		}
	}
	if protocol == "smtp" || runAll {
		if judge.TcpSMTP(result) {
			printSuccess("TCP/SMTP", result)
			return true
		}
	}
	if protocol == "imap" || runAll {
		if judge.TcpIMAP(result) {
			printSuccess("TCP/IMAP", result)
			return true
		}
	}
	if protocol == "ssh" || runAll {
		if judge.TcpSSH(result) {
			printSuccess("TCP/SSH", result)
			return true
		}
	}
	if protocol == "pop3" || runAll {
		if judge.TcpPOP3(result) {
			printSuccess("TCP/POP3", result)
			return true
		}
	}
	if protocol == "vnc" || runAll {
		if judge.TcpVNC(result) {
			printSuccess("TCP/VNC", result)
			return true
		}
	}
	if protocol == "telnet" || runAll {
		if judge.TcpTelnet(result) {
			printSuccess("TCP/Telnet", result)
			return true
		}
	}
	if protocol == "ftp" || runAll {
		if judge.TcpFTP(result) {
			printSuccess("TCP/FTP", result)
			return true
		}
	}
	if protocol == "snmp" || runAll {
		if judge.TcpSNMP(result) {
			printSuccess("TCP/SNMP", result)
			return true
		}
	}
	if protocol == "oracle" || runAll {
		if judge.TcpOracle(result, Args) {
			printSuccess("TCP/Oracle", result)
			return true
		}
	}
	if protocol == "frp" || runAll {
		if judge.TcpFrp(result, Args) {
			printSuccess("TCP/Frp", result)
			return true
		}
	}
	if protocol == "socks" || runAll {
		if judge.TcpSocks(result, Args) {
			printSuccess("TCP/Socks", result)
			return true
		}
	}
	if protocol == "ldap" || runAll {
		if judge.TcpLDAP(result, Args) {
			printSuccess("TCP/LDAP", result)
			return true
		}
	}
	if protocol == "rmi" || runAll {
		if judge.TcpRMI(result, Args) {
			printSuccess("TCP/RMI", result)
			return true
		}
	}
	if protocol == "activemq" || runAll {
		if judge.TcpActiveMQ(result) {
			printSuccess("TCP/ActiveMQ", result)
			return true
		}
	}
	if protocol == "rtsp" || runAll {
		if judge.TcpRTSP(result, Args) {
			printSuccess("TCP/RTSP", result)
			return true
		}
	}
	if protocol == "rdp" || runAll {
		if judge.TcpRDP(result, Args) {
			printSuccess("TCP/RDP", result)
			return true
		}
	}

	if protocol == "dcerpc" || runAll {
		if judge.TcpDceRpc(result, Args) {
			printSuccess("TCP/DceRpc", result)
			return true
		}
	}
	if protocol == "mssql" || runAll {
		if judge.TcpMssql(result, Args) {
			printSuccess("TCP/Mssql", result)
			return true
		}
	}
	if protocol == "smb" || runAll {
		if judge.TcpSMB(result, Args) {
			printSuccess("TCP/SMB", result)
			return true
		}
	}

	status := result["status"].(string)
	if status == "open" && runAll {
		printFailed("TCP/unknown", result)
	}
	return false
}

func JudgeTls(result map[string]interface{}, Args map[string]interface{}) bool {
	protocol := result["protocol"].(string)
	runAll := true
	if protocol != "" {
		runAll = false
	}
	if protocol == "http" || protocol == "https" || runAll {
		if judge.TlsHTTPS(result, Args) {
			printSuccess("TLS/HTTPS", result)
			return true
		}
	}
	if protocol == "rdp" || runAll {
		if judge.TlsRDP(result, Args) {
			printSuccess("TLS/RDP", result)
			return true
		}
	}
	if protocol == "redis-ssl" || runAll {
		if judge.TlsRedisSsl(result) {
			printSuccess("TLS/Redis-ssl", result)
			return true
		}
	}

	status := result["status"].(string)
	if status == "open" && runAll {
		printFailed("TLS/unknown", result)
	}
	return false
}

func JudgeUdp(result map[string]interface{}, Args map[string]interface{}) bool {
	protocol := result["protocol"].(string)
	runAll := true
	if protocol != "" {
		runAll = false
	}
	if protocol == "nbns" || runAll {
		if judge.UdpNbns(result, Args) {
			printSuccess("UDP/NBNS", result)
			return true
		}
	}

	var buffer [256]byte
	status := result["status"].(string)
	if bytes.Equal(result["banner.byte"].([]byte), buffer[:]) {
		result["status"] = "close"
		return false
	} else if status == "open" && runAll {
		printFailed("UDP/unknown", result)
		//logger.Failed(fmt.Sprintf("[%s] %s [%s]",logger.Cyan("UDP/unknown"), parse.SchemeParse(result), logger.Blue(banner)))
	}
	return false
}

func printSuccess(protocol string, result map[string]interface{}) {
	success, b := result["identify.bool"].(bool)
	if b == false {
		logger.Success(fmt.Sprintf("[%s] %s [%s]",
			logger.Cyan(protocol),
			parse.SchemeParse(result),
			logger.Blue(result["banner.string"].(string))),
		)
		result["identify.string"] = logger.Clean(result["identify.string"].(string))
		result["note"] = logger.Clean(result["note"].(string))
		return
	}

	if success {
		logger.Success(fmt.Sprintf("[%s] %s %s [%s]",
			logger.Cyan(protocol),
			result["identify.string"].(string),
			parse.SchemeParse(result),
			logger.Blue(result["banner.string"].(string))),
		)
		result["identify.string"] = logger.Clean(result["identify.string"].(string))
		result["note"] = logger.Clean(result["note"].(string))
		return
	} else {
		logger.Success(fmt.Sprintf("[%s] %s [%s]",
			logger.Cyan(protocol),
			parse.SchemeParse(result),
			logger.Blue(result["banner.string"].(string))),
		)
		result["identify.string"] = logger.Clean(result["identify.string"].(string))
		result["note"] = logger.Clean(result["note"].(string))
		return
	}
}

func printFailed(p string, result map[string]interface{}) {
	if result["status"].(string) == "open" {
		logger.Failed(fmt.Sprintf("[%s] %s [%s]",
			logger.Cyan(p),
			parse.SchemeParse(result),
			logger.Blue(result["banner.string"].(string))),
		)
		result["identify.string"] = logger.Clean(result["identify.string"].(string))
		result["note"] = logger.Clean(result["note"].(string))
	}
}
