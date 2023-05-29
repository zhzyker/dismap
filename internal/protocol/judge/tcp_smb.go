package judge

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/proxy"

	crand "crypto/rand"
)

func TcpSMB(result *model.Result) bool {
	if smbPublic(result, result.Host, result.Port, flag.Timeout) {
		return true
	}
	res, data, err := smb2(result.Host, result.Port, flag.Timeout)
	if err != nil || res["ntlmssp.Version"] == "" {
		return false
	} else {
		banner := fmt.Sprintf("Version:%s||DNSComputer:%s||TargetName:%s||NetbiosComputer:%s",
			res["ntlmssp.Version"],
			res["ntlmssp.DNSComputer"],
			res["ntlmssp.TargetName"],
			res["ntlmssp.NetbiosComputer"],
		)
		result.Banner = banner
		result.Protocol = "smb"
		result.BannerB = data
		return true
	}
}

func smbPublic(result *model.Result, host string, port int, timeout int) bool {
	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return false
	}

	msg1 := "\x00\x00\x00\x85\xff\x53\x4d\x42\x72\x00\x00\x00\x00\x18\x53\xc0\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xff\xfe\x00\x00\x40\x00\x00\x62\x00\x02\x50\x43\x20\x4e\x45\x54\x57\x4f\x52\x4b\x20\x50\x52\x4f\x47\x52\x41\x4d\x20\x31\x2e\x30\x00\x02\x4c\x41\x4e\x4d\x41\x4e\x31\x2e\x30\x00\x02\x57\x69\x6e\x64\x6f\x77\x73\x20\x66\x6f\x72\x20\x57\x6f\x72\x6b\x67\x72\x6f\x75\x70\x73\x20\x33\x2e\x31\x61\x00\x02\x4c\x4d\x31\x2e\x32\x58\x30\x30\x32\x00\x02\x4c\x41\x4e\x4d\x41\x4e\x32\x2e\x31\x00\x02\x4e\x54\x20\x4c\x4d\x20\x30\x2e\x31\x32\x00"
	msg2 := "\x00\x00\x00\x88\xff\x53\x4d\x42\x73\x00\x00\x00\x00\x18\x07\xc0\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xff\xfe\x00\x00\x40\x00\x0d\xff\x00\x88\x00\x04\x11\x0a\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\xd4\x00\x00\x00\x4b\x00\x00\x00\x00\x00\x00\x57\x00\x69\x00\x6e\x00\x64\x00\x6f\x00\x77\x00\x73\x00\x20\x00\x32\x00\x30\x00\x30\x00\x30\x00\x20\x00\x32\x00\x31\x00\x39\x00\x35\x00\x00\x00\x57\x00\x69\x00\x6e\x00\x64\x00\x6f\x00\x77\x00\x73\x00\x20\x00\x32\x00\x30\x00\x30\x00\x30\x00\x20\x00\x35\x00\x2e\x00\x30\x00\x00\x00"
	_, err = conn.Write([]byte(msg1))
	if err != nil {
		return false
	}
	reply1 := make([]byte, 256)
	_, _ = conn.Read(reply1)

	if hex.EncodeToString(reply1[0:8]) != "00000081ff534d42" {
		return false
	}
	_, err = conn.Write([]byte(msg2))
	if err != nil {
		return false
	}
	reply2 := make([]byte, 512)
	_, _ = conn.Read(reply2)
	if conn != nil {
		_ = conn.Close()
	}

	var buffer bytes.Buffer
	for i := 0; i < len(reply2[46:]); {
		b := reply2[46:][i : i+2]
		i += 2
		if 46+i == len(reply2[46:]) {
			break
		}
		if string(b) == "\x00\x00" {
			if string(reply2[46+i+2:46+i+2+2]) == "\x00\x00" {
				break
			}
			buffer.Write([]byte("\x7C\x7C"))
			result.Banner = strings.Join([]string{buffer.String()}, ",")
			continue
		}
		buffer.Write(b[0:1])
		result.Banner = strings.Join([]string{buffer.String()}, ",")
	}

	var ban [512]byte
	if bytes.Equal(reply2[:], ban[:]) {
		return false
	} else {
		result.Protocol = "smb"
		result.BannerB = reply2
		return true
	}
}

// smb2 from https://github.com/RumbleDiscovery/rumble-tools/blob/main/cmd/rumble-smb2-sessions/main.go
func smb2(host string, port int, timeout int) (map[string]string, []byte, error) {
	// SMB1NegotiateProtocolRequest is a SMB1 request that advertises support for SMB2
	var smb1NegotiateProtocolRequest = []byte{
		0x00, 0x00, 0x00, 0xd4, 0xff, 0x53, 0x4d, 0x42,
		0x72, 0x00, 0x00, 0x00, 0x00, 0x18, 0x43, 0xc8,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xfe, 0xff,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xb1, 0x00, 0x02,
		0x50, 0x43, 0x20, 0x4e, 0x45, 0x54, 0x57, 0x4f,
		0x52, 0x4b, 0x20, 0x50, 0x52, 0x4f, 0x47, 0x52,
		0x41, 0x4d, 0x20, 0x31, 0x2e, 0x30, 0x00, 0x02,
		0x4d, 0x49, 0x43, 0x52, 0x4f, 0x53, 0x4f, 0x46,
		0x54, 0x20, 0x4e, 0x45, 0x54, 0x57, 0x4f, 0x52,
		0x4b, 0x53, 0x20, 0x31, 0x2e, 0x30, 0x33, 0x00,
		0x02, 0x4d, 0x49, 0x43, 0x52, 0x4f, 0x53, 0x4f,
		0x46, 0x54, 0x20, 0x4e, 0x45, 0x54, 0x57, 0x4f,
		0x52, 0x4b, 0x53, 0x20, 0x33, 0x2e, 0x30, 0x00,
		0x02, 0x4c, 0x41, 0x4e, 0x4d, 0x41, 0x4e, 0x31,
		0x2e, 0x30, 0x00, 0x02, 0x4c, 0x4d, 0x31, 0x2e,
		0x32, 0x58, 0x30, 0x30, 0x32, 0x00, 0x02, 0x44,
		0x4f, 0x53, 0x20, 0x4c, 0x41, 0x4e, 0x4d, 0x41,
		0x4e, 0x32, 0x2e, 0x31, 0x00, 0x02, 0x4c, 0x41,
		0x4e, 0x4d, 0x41, 0x4e, 0x32, 0x2e, 0x31, 0x00,
		0x02, 0x53, 0x61, 0x6d, 0x62, 0x61, 0x00, 0x02,
		0x4e, 0x54, 0x20, 0x4c, 0x41, 0x4e, 0x4d, 0x41,
		0x4e, 0x20, 0x31, 0x2e, 0x30, 0x00, 0x02, 0x4e,
		0x54, 0x20, 0x4c, 0x4d, 0x20, 0x30, 0x2e, 0x31,
		0x32, 0x00, 0x02, 0x53, 0x4d, 0x42, 0x20, 0x32,
		0x2e, 0x30, 0x30, 0x32, 0x00, 0x02, 0x53, 0x4d,
		0x42, 0x20, 0x32, 0x2e, 0x3f, 0x3f, 0x3f, 0x00,
	}

	info := make(map[string]string)
	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if err != nil {
		return info, []byte{}, err
	}
	err = SMBSendData(conn, smb1NegotiateProtocolRequest, timeout)
	if err != nil {
		return info, []byte{}, err
	}

	_, err = SMBReadFrame(conn, timeout)
	if err != nil {
		return info, []byte{}, err
	}

	err = SMBSendData(conn, SMB2NegotiateProtocolRequest(host), timeout)
	if err != nil {
		return info, []byte{}, err
	}

	data, _ := SMBReadFrame(conn, timeout)
	SMB2ExtractFieldsFromNegotiateReply(data, info)

	// SMB2SessionSetupNTLMSSP is a SMB2 SessionSetup NTLMSSP request
	var smb2SessionSetupNTLMSSP = []byte{
		0x00, 0x00, 0x00, 0xa2, 0xfe, 0x53, 0x4d, 0x42,
		0x40, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x19, 0x00, 0x00, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x58, 0x00, 0x4a, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x60, 0x48, 0x06, 0x06,
		0x2b, 0x06, 0x01, 0x05, 0x05, 0x02, 0xa0, 0x3e,
		0x30, 0x3c, 0xa0, 0x0e, 0x30, 0x0c, 0x06, 0x0a,
		0x2b, 0x06, 0x01, 0x04, 0x01, 0x82, 0x37, 0x02,
		0x02, 0x0a, 0xa2, 0x2a, 0x04, 0x28, 0x4e, 0x54,
		0x4c, 0x4d, 0x53, 0x53, 0x50, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x97, 0x82, 0x08, 0xe2, 0x00, 0x00,
		0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x0a, 0x00,
		0xba, 0x47, 0x00, 0x00, 0x00, 0x0f,
	}

	setup := make([]byte, len(smb2SessionSetupNTLMSSP))
	copy(setup, smb2SessionSetupNTLMSSP)

	// Set the ProcessID
	binary.LittleEndian.PutUint16(setup[4+32:], 0xfeff)

	err = SMBSendData(conn, setup, timeout)
	if err != nil {
		return info, []byte{}, err
	}

	data, err = SMBReadFrame(conn, timeout)
	SMB2ExtractSIDFromSessionSetupReply(data, info)
	SMBExtractFieldsFromSecurityBlob(data, info)

	return info, data, err
}

// RandomBytes generates a random byte sequence of the requested length
func RandomBytes(numBytes int) []byte {
	randBytes := make([]byte, numBytes)
	// err := binary.Read(crand.Reader, binary.BigEndian, &randBytes)
	binary.Read(crand.Reader, binary.BigEndian, &randBytes)
	return randBytes
}

// SMBReadFrame reads the netBios header then the full response
func SMBReadFrame(conn net.Conn, t int) ([]byte, error) {
	timeout := time.Now().Add(time.Duration(t) * time.Second)
	res := []byte{}
	nbh := make([]byte, 4)

	err := conn.SetReadDeadline(timeout)
	if err != nil {
		return res, err
	}

	// Read the NetBIOS header
	n, err := conn.Read(nbh[:])
	if err != nil {
		// Return if EOF is reached
		if err == io.EOF {
			return res, nil
		}

		// Return if timeout is reached
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return res, nil
		}

		// If we have data and received an error, it was probably a reset
		if len(res) > 0 {
			return res, nil
		}

		return res, err
	}

	if n != 4 {
		return res, nil
	}

	res = append(res[:], nbh[:n]...)
	dlen := binary.BigEndian.Uint32(nbh[:]) & 0x00ffffff
	buf := make([]byte, dlen)
	n, err = conn.Read(buf[:])
	if err != nil {
		// Return if EOF is reached
		if err == io.EOF {
			return res, nil
		}

		// Return if timeout is reached
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return res, nil
		}

		// If we have data and received an error, it was probably a reset
		if len(res) > 0 {
			return res, nil
		}

		return res, err
	}
	res = append(res[:], buf[:n]...)
	return res, nil
}

// SMB2ExtractFieldsFromNegotiateReply extracts useful fields from the SMB2 negotiate response
func SMB2ExtractFieldsFromNegotiateReply(blob []byte, info map[string]string) {

	smbOffset := bytes.Index(blob, []byte{0xfe, 'S', 'M', 'B'})
	if smbOffset < 0 {
		return
	}

	data := blob[smbOffset:]

	// Basic sanity check
	if len(data) < (64 + 8 + 16 + 36) {
		return
	}

	switch binary.LittleEndian.Uint16(data[64+2:]) {
	case 0:
		info["smb.Signing"] = "disabled"
	case 1:
		info["smb.Signing"] = "enabled"
	case 2, 3:
		info["smb.Signing"] = "required"
	}

	info["smb.Dialect"] = fmt.Sprintf("0x%.4x", binary.LittleEndian.Uint16(data[64+4:]))
	//info["smb.GUID"] = uuid.FromBytesOrNil(data[64+8 : 64+8+16]).String()
	info["smb.Capabilities"] = fmt.Sprintf("0x%.8x", binary.LittleEndian.Uint32(data[64+8+16:]))

	negCtxCount := int(binary.LittleEndian.Uint16(data[64+6:]))
	negCtxOffset := int(binary.LittleEndian.Uint32(data[64+8+16+36:]))
	if negCtxCount == 0 || negCtxOffset == 0 || negCtxOffset+(negCtxCount*8) > len(data) {
		return
	}

	negCtxData := data[negCtxOffset:]
	idx := 0
	for {
		if idx+8 > len(negCtxData) {
			break
		}
		negType := int(binary.LittleEndian.Uint16(negCtxData[idx:]))
		negLen := int(binary.LittleEndian.Uint16(negCtxData[idx+2:]))
		idx += 8

		if idx+negLen > len(negCtxData) {
			break
		}
		negData := negCtxData[idx : idx+negLen]

		SMB2ParseNegotiateContext(negType, negData, info)

		// Move the index to the next context
		idx += negLen
		// Negotiate Contexts are aligned on 64-bit boundaries
		for idx%8 != 0 {
			idx++
		}
	}
}

// SMB2ParseNegotiateContext decodes fields from the SMB2 Negotiate Context values
func SMB2ParseNegotiateContext(t int, data []byte, info map[string]string) {
	switch t {
	case 1:
		// SMB2_PREAUTH_INTEGRITY_CAPABILITIES
		if len(data) < 6 {
			return
		}
		hashCount := int(binary.LittleEndian.Uint16(data[:]))
		// MUST only be one in responses
		if hashCount != 1 {
			return
		}
		hashSaltLen := int(binary.LittleEndian.Uint16(data[2:]))
		hashType := int(binary.LittleEndian.Uint16(data[4:]))
		hashName := "sha512"
		if hashType != 1 {
			hashName = fmt.Sprintf("unknown-%d", hashType)
		}
		info["smb.HashAlg"] = hashName
		info["smb.HashSaltLen"] = fmt.Sprintf("%d", hashSaltLen)

	case 2:
		// SMB2_ENCRYPTION_CAPABILITIES
		if len(data) < 4 {
			return
		}
		cipherCount := int(binary.LittleEndian.Uint16(data[:]))
		if len(data) < 2+(2*cipherCount) {
			return
		}
		// MUST only be one in responses
		if cipherCount != 1 {
			return
		}
		cipherList := []string{}
		for i := 0; i < cipherCount; i++ {
			cipherID := int(binary.LittleEndian.Uint16(data[2+(i*2):]))
			cipherName := ""
			switch cipherID {
			case 1:
				cipherName = "aes-128-ccm"
			case 2:
				cipherName = "aes-128-gcm"
			default:
				cipherName = fmt.Sprintf("unknown-%d", cipherID)
			}
			cipherList = append(cipherList, cipherName)
		}

		info["smb.CipherAlg"] = strings.Join(cipherList, "\t")

	case 3:
		// SMB2_COMPRESSION_CAPABILITIES
		if len(data) < 10 {
			return
		}
		compCount := int(binary.LittleEndian.Uint16(data[:]))
		if len(data) < 2+2+4+(2*compCount) {
			return
		}
		// MUST only be one in responses
		if compCount != 1 {
			return
		}
		compList := []string{}
		for i := 0; i < compCount; i++ {
			compID := int(binary.LittleEndian.Uint16(data[8+(i*2):]))
			compName := ""
			switch compID {
			case 0:
				compName = "none"
			case 1:
				compName = "lznt1"
			case 2:
				compName = "lz77"
			case 3:
				compName = "lz77+huff"
			case 4:
				compName = "patternv1"
			default:
				compName = fmt.Sprintf("unknown-%d", compID)
			}
			compList = append(compList, compName)
		}
		info["smb.CompressionFlags"] = fmt.Sprintf("0x%.4x", binary.LittleEndian.Uint32(data[4:]))
		info["smb.CompressionAlg"] = strings.Join(compList, "\t")
	}
}

// SMBSendData writes a SMB request to a socket
func SMBSendData(conn net.Conn, data []byte, timeout int) error {
	err := conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	if err != nil {
		return err
	}

	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	_ = n

	return nil
}

// SMB2ExtractSIDFromSessionSetupReply tries to extract the SessionID and Signature from a SMB2 reply
func SMB2ExtractSIDFromSessionSetupReply(blob []byte, info map[string]string) {
	smbOffset := bytes.Index(blob, []byte{0xfe, 'S', 'M', 'B'})
	if smbOffset < 0 {
		return
	}

	smbData := blob[smbOffset:]
	if len(smbData) < 48 {
		return
	}

	status := binary.LittleEndian.Uint32(smbData[8:])
	info["smb.Status"] = fmt.Sprintf("0x%.8x", status)

	sessID := binary.LittleEndian.Uint64(smbData[40:])
	info["smb.SessionID"] = fmt.Sprintf("0x%.16x", sessID)

	if len(smbData) >= 64 {
		sigData := hex.EncodeToString(smbData[48:64])
		if sigData != "00000000000000000000000000000000" {
			info["smb.Signature"] = sigData
		}
	}
}

// SMBExtractValueFromOffset peels a field out of a SMB buffer
func SMBExtractValueFromOffset(blob []byte, idx int) ([]byte, int, error) {
	res := []byte{}

	if len(blob) < (idx + 6) {
		return res, idx, fmt.Errorf("data truncated")
	}

	len1 := binary.LittleEndian.Uint16(blob[idx:])
	idx += 2

	// len2 := binary.LittleEndian.Uint16(blob[idx:])
	idx += 2

	off := binary.LittleEndian.Uint32(blob[idx:])
	idx += 4

	// Allow zero length values
	if len1 == 0 {
		return res, idx, nil
	}

	if len(blob) < int(off+uint32(len1)) {
		return res, idx, fmt.Errorf("data value truncated")
	}

	res = append(res, blob[off:off+uint32(len1)]...)
	return res, idx, nil
}

// SMBExtractFieldsFromSecurityBlob extracts fields from the NTLMSSP response
func SMBExtractFieldsFromSecurityBlob(blob []byte, info map[string]string) {
	var err error

	ntlmsspOffset := bytes.Index(blob, []byte{'N', 'T', 'L', 'M', 'S', 'S', 'P', 0x00, 0x02, 0x00, 0x00, 0x00})
	if ntlmsspOffset < 0 {
		return
	}

	data := blob[ntlmsspOffset:]

	// Basic sanity check
	if len(data) < (12 + 6 + 12 + 8 + 6 + 8) {
		return
	}

	idx := 12

	targetName, idx, err := SMBExtractValueFromOffset(data, idx)
	if err != nil {
		return
	}

	// Negotiate Flags
	negotiateFlags := binary.LittleEndian.Uint32(data[idx:])
	info["ntlmssp.NegotiationFlags"] = fmt.Sprintf("0x%.8x", negotiateFlags)
	idx += 4

	// NTLM Server Challenge
	idx += 8

	// Reserved
	idx += 8

	// Target Info
	targetInfo, idx, err := SMBExtractValueFromOffset(data, idx)
	if err != nil {
		return
	}

	// Version
	versionMajor := uint8(data[idx])
	idx++

	versionMinor := uint8(data[idx])
	idx++

	versionBuild := binary.LittleEndian.Uint16(data[idx:])
	idx += 2

	ntlmRevision := binary.BigEndian.Uint32(data[idx:])

	// macOS reverses the endian order of this field for some reason
	if ntlmRevision == 251658240 {
		ntlmRevision = binary.LittleEndian.Uint32(data[idx:])
	}

	info["ntlmssp.Version"] = fmt.Sprintf("%d.%d.%d", versionMajor, versionMinor, versionBuild)
	info["ntlmssp.NTLMRevision"] = fmt.Sprintf("%d", ntlmRevision)
	info["ntlmssp.TargetName"] = TrimName(string(targetName))

	idx = 0
	for {
		if idx+4 > len(targetInfo) {
			break
		}

		attrType := binary.LittleEndian.Uint16(targetInfo[idx:])
		idx += 2

		// End of List
		if attrType == 0 {
			break
		}

		attrLen := binary.LittleEndian.Uint16(targetInfo[idx:])
		idx += 2

		if idx+int(attrLen) > len(targetInfo) {
			// log.Printf("too short: %d/%d", idx+int(attrLen), len(targetInfo))
			break
		}

		attrVal := targetInfo[idx : idx+int(attrLen)]
		idx += int(attrLen)

		switch attrType {
		case 1:
			info["ntlmssp.NetbiosComputer"] = TrimName(string(attrVal))
		case 2:
			info["ntlmssp.NetbiosDomain"] = TrimName(string(attrVal))
		case 3:
			info["ntlmssp.DNSComputer"] = TrimName(string(attrVal))
		case 4:
			info["ntlmssp.DNSDomain"] = TrimName(string(attrVal))
		case 7:
			ts := binary.LittleEndian.Uint64(attrVal[:])
			info["ntlmssp.Timestamp"] = fmt.Sprintf("0x%.16x", ts)

		}

		// End of List
		if attrType == 0 {
			break
		}
	}
}

// TrimName removes null bytes and trims leading and trailing spaces from a string
func TrimName(name string) string {
	return strings.TrimSpace(strings.Replace(name, "\x00", "", -1))
}

// SMB2NegotiateProtocolRequest generates a new Negotiate request with the specified target name
func SMB2NegotiateProtocolRequest(dst string) []byte {

	base := []byte{
		0xfe, 0x53, 0x4d, 0x42,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xff, 0xfe, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x24, 0x00, 0x05, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x7f, 0x00, 0x00, 0x00,
	}

	// Client GUID (16)
	base = append(base[:], RandomBytes(16)...)

	base = append(base[:], []byte{
		0x70, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
		0x02, 0x02, 0x10, 0x02, 0x00, 0x03, 0x02, 0x03,
		0x11, 0x03, 0x00, 0x00, 0x01, 0x00, 0x26, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x20, 0x00,
		0x01, 0x00,
	}...)

	// SHA-512 Salt (32)
	base = append(base[:], RandomBytes(32)...)

	base = append(base[:], []byte{
		0x00, 0x00, 0x02, 0x00, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x02, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x03, 0x00, 0x0e, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x03, 0x00,
		0x01, 0x00, 0x00, 0x00,
	}...)

	encodedDst := make([]byte, len(dst)*2)
	for i, b := range []byte(dst) {
		encodedDst[i*2] = b
	}

	netname := []byte{0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	binary.LittleEndian.PutUint16(netname[2:], uint16(len(encodedDst)))
	netname = append(netname, encodedDst...)

	base = append(base, netname...)

	nbhd := make([]byte, 4)
	binary.BigEndian.PutUint32(nbhd, uint32(len(base)))
	nbhd = append(nbhd, base...)
	return nbhd
}
