package http

import (
	"fmt"
	"net/http"
	"strings"
)

// getHeaders 函数用于从 http.Response 中读取响应头信息并返回字符串形式的响应头内容
func getHeaders(response *http.Response) string {
	// 首先判断响应是否为空或者响应头是否为空，如果是则返回空字符串
	if response == nil || response.Header == nil {
		return ""
	}
	var builder strings.Builder
	// 遍历响应头，将每个头字段名和对应值拼接成字符串并加入到 builder 中
	for name, values := range response.Header {
		builder.WriteString(fmt.Sprintf("%s: %s\n", name, strings.Join(values, ", ")))
	}
	//logger.DBG("Target url title: \n" + builder.String())
	return builder.String()
}
