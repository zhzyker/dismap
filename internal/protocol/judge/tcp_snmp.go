package judge

import (
	"bytes"
)

func TcpSNMP(result map[string]interface{}) bool {
	b := result["banner.byte"].([]byte)
	if bytes.Equal(b[:], make([]byte, 0)[:]) {
		return false
	}

	buff := []byte{0x41, 0x01, 0x02}
	snmp := result["banner.byte"].([]byte)[0:3]
	if bytes.Equal(buff[:], snmp[:]) {
		result["protocol"] = "snmp"
		return true
	}
	return false
}