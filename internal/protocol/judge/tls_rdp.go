package judge

func TlsRDP(result map[string]interface{}, Args map[string]interface{}) bool {
	if TcpRDP(result, Args) {
		return true
	}
	return false
}