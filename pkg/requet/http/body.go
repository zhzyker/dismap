package http

import (
	"bytes"
	"io"
	"net/http"
)

// getBody 函数用于从 http.Response 中读取响应体并返回字符串形式的响应体内容
func getBody(response *http.Response) string {
	// 首先判断响应是否为空或者 body 是否为空，如果是则返回空字符串
	if response == nil || response.Header == nil {
		return ""
	}
	var buf bytes.Buffer // 创建一个 bytes.Buffer 对象
	if _, err := io.Copy(&buf, response.Body); err != nil {
		return ""
	}
	// logger.DBG("Target url body: \n" + buf.String())
	return buf.String()
}
