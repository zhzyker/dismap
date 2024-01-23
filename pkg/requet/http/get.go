package http

// HttpResult 结构体包含了 HTTP 请求结果的各个字段
type HttpResult struct {
	Url        string
	Body       string
	Headers    string
	StatusCode int
	Title      string
}

// Get 获取 HTTP 的响应数据并返回到 httpResult , 用于其他方法调用 HTTP 请求结果
func Get(url string) (HttpResult, error) {
	client := newHTTPClient()
	resp, err := client.Get(url)
	var code int
	defer func() {
		if err != nil {
			code = 0
		} else if resp != nil {
			// 获取 HTTP 状态码
			code = resp.StatusCode
			if err := resp.Body.Close(); err != nil {
			}
		}
	}()
	// 获取 HTTP 请求结果的 Body
	body := getBody(resp)
	return HttpResult{
		url,
		body,
		getHeaders(resp),
		code,
		getTitle(body),
	}, err
}
