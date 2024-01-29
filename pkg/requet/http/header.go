package http

import (
	"fmt"
	"net/http"
	"strings"
)

// getHeader 函数用于从 http.Response 中读取响应头信息并返回字符串形式的响应头内容
func getHeader(req *http.Response) []byte {
	var headers []byte
	// 遍历响应头，将每个头字段名和对应值按照格式 (key: value) 拼接成字符串并加入到 header 中
	for name, values := range req.Header {
		headers = append(headers, []byte(fmt.Sprintf("%s: %s\n", name, strings.Join(values, ", ")))...)
	}
	return headers
}
