package http

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

// newHTTPClient 返回一个新的 *http.Client，其中包含以下配置:
// - Timeout: 设置HTTP客户端的超时时间为 10 秒
// - Transport: 设置HTTP客户端的传输层配置，其中包括:
// - TLSClientConfig: 配置TLS客户端的验证方式为跳过证书验证
// - DisableCompression: 禁用 Accept-Encoding: gzip, deflate 压缩,若开启压缩,无法得到正常的响应 Body 长度,导致分块读取 Body 不理想
// - CheckRedirect: 允许重定向
func newHTTPClient(timeout int, proxy *url.URL) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			DisableCompression: true,
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			Proxy:              http.ProxyURL(proxy),
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
}
