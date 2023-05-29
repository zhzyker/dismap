package judge

import (
	"bytes"

	"github.com/zhzyker/dismap/internal/model"
)

func TcpSNMP(result *model.Result) bool {
	b := result.BannerB
	if bytes.Equal(b[:], make([]byte, 0)[:]) {
		return false
	}

	buff := []byte{0x41, 0x01, 0x02}
	snmp := result.BannerB[0:3]
	if bytes.Equal(buff[:], snmp[:]) {
		result.Protocol = "snmp"
		return true
	}
	return false
}
