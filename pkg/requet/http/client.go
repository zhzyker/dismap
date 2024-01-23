package http

import (
	"crypto/tls"
	"net/http"
	"time"
)

// newHTTPClient 返回一个新的 *http.Client，其中包含以下配置:
// - Timeout: 设置HTTP客户端的超时时间为 10 秒
// - Transport: 设置HTTP客户端的传输层配置，其中包括:
//   - TLSClientConfig: 配置TLS客户端的验证方式为跳过证书验证
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(10) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
