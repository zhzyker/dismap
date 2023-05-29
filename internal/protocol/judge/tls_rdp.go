package judge

import "github.com/zhzyker/dismap/internal/model"

func TlsRDP(result *model.Result) bool {
	if TcpRDP(result) {
		return true
	}
	return false
}
