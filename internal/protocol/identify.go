package protocol

import (
	"bytes"
	"fmt"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/protocol/judge"
	"github.com/zhzyker/dismap/pkg/logger"
)

func JudgeTcp(result *model.Result) bool {
	protocol := result.Protocol
	runAll := true
	if protocol != "" {
		runAll = false
	}
	if protocol == "http" || protocol == "https" || runAll {
		if judge.TcpHTTP(result) {
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
		if judge.TcpOracle(result) {
			printSuccess("TCP/Oracle", result)
			return true
		}
	}
	if protocol == "frp" || runAll {
		if judge.TcpFrp(result) {
			printSuccess("TCP/Frp", result)
			return true
		}
	}
	if protocol == "socks" || runAll {
		if judge.TcpSocks(result) {
			printSuccess("TCP/Socks", result)
			return true
		}
	}
	if protocol == "ldap" || runAll {
		if judge.TcpLDAP(result) {
			printSuccess("TCP/LDAP", result)
			return true
		}
	}
	if protocol == "rmi" || runAll {
		if judge.TcpRMI(result) {
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
		if judge.TcpRTSP(result) {
			printSuccess("TCP/RTSP", result)
			return true
		}
	}
	if protocol == "rdp" || runAll {
		if judge.TcpRDP(result) {
			printSuccess("TCP/RDP", result)
			return true
		}
	}

	if protocol == "dcerpc" || runAll {
		if judge.TcpDceRpc(result) {
			printSuccess("TCP/DceRpc", result)
			return true
		}
	}
	if protocol == "mssql" || runAll {
		if judge.TcpMssql(result) {
			printSuccess("TCP/Mssql", result)
			return true
		}
	}
	if protocol == "smb" || runAll {
		if judge.TcpSMB(result) {
			printSuccess("TCP/SMB", result)
			return true
		}
	}
	if protocol == "giop" || runAll {
		if judge.TcpGIOP(result) {
			printSuccess("TCP/GIOP", result)
			return true
		}
	}

	status := result.Status
	if status == "open" && runAll {
		printFailed("TCP/unknown", result)
	}
	return false
}

func JudgeTls(result *model.Result) bool {
	protocol := result.Protocol
	runAll := true
	if protocol != "" {
		runAll = false
	}
	if protocol == "http" || protocol == "https" || runAll {
		if judge.TlsHTTPS(result) {
			printSuccess("TLS/HTTPS", result)
			return true
		}
	}
	if protocol == "rdp" || runAll {
		if judge.TlsRDP(result) {
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

	status := result.Status
	if status == "open" && runAll {
		printFailed("TLS/unknown", result)
	}
	return false
}

func JudgeUdp(result *model.Result) bool {
	protocol := result.Protocol
	runAll := true
	if protocol != "" {
		runAll = false
	}
	if protocol == "nbns" || runAll {
		if judge.UdpNbns(result) {
			printSuccess("UDP/NBNS", result)
			return true
		}
	}

	var buffer [256]byte
	status := result.Status
	if bytes.Equal(result.BannerB, buffer[:]) {
		result.Status = "close"
		return false
	} else if status == "open" && runAll {
		printFailed("UDP/unknown", result)
		//logger.Failed(fmt.Sprintf("[%s] %s [%s]",logger.Cyan("UDP/unknown"), parse.SchemeParse(result), logger.Blue(banner)))
	}
	return false
}

func printSuccess(protocol string, result *model.Result) {
	success := result.IdentifyBool

	if success {
		logger.Success(fmt.Sprintf("[%s] %s %s [%s]",
			logger.Cyan(protocol),
			result.IdentifyStr,
			parse.SchemeParse(result),
			logger.Blue(result.Banner)),
		)
		result.IdentifyStr = logger.Clean(result.IdentifyStr)
		result.Note = logger.Clean(result.Note)
		return
	} else {
		logger.Success(fmt.Sprintf("[%s] %s [%s]",
			logger.Cyan(protocol),
			parse.SchemeParse(result),
			logger.Blue(result.Banner)),
		)
		result.IdentifyStr = logger.Clean(result.IdentifyStr)
		result.Note = logger.Clean(result.Note)
		return
	}
}

func printFailed(p string, result *model.Result) {
	if result.Status == "open" {
		logger.Failed(fmt.Sprintf("[%s] %s [%s]",
			logger.Cyan(p),
			parse.SchemeParse(result),
			logger.Blue(result.Banner)),
		)
		result.IdentifyStr = logger.Clean(result.IdentifyStr)
		result.Note = logger.Clean(result.Note)
	}
}
