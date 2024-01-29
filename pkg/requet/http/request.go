package http

import "net/url"

// Responses 结构体包含了 HTTP 请求结果的各个字段
type Responses struct {
	Url        string
	Body       []byte
	Header     []byte
	StatusCode int
	Title      string
	Favicon    []byte
}

// Request 获取 HTTP 的响应数据并返回到 Responses ,包括 Url Body Headers StatusCode Title ,用于其他方法调用 HTTP 请求结果.
func Request(url string, timeout int, isFavicon bool, proxy *url.URL) (Responses, error) {
	client := newHTTPClient(timeout, proxy)
	req, err := client.Get(url)
	if err != nil {
		return Responses{}, err // 当 HTTP 请求发送失败时,返回错误信息
	}
	defer func() {
		if req != nil {
			if err := req.Body.Close(); err != nil {
				// 处理关闭 Body 过程中的错误,可以记录日志等
			}
		}
	}()
	body := getBody(req)
	return Responses{
		url,
		body,
		getHeader(req),
		req.StatusCode,
		getTitle(body),
		getFavicon(url, client, isFavicon),
	}, nil
}
