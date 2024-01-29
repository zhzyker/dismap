package http

import (
	"bytes"
	"io"
	"net/http"
)

// getBody 函数用于从 http.Response 中读取响应体并返回字符串形式的响应体内容
func getBody(req *http.Response) []byte {
	var buf bytes.Buffer
	// 判断响应体大小, 如果小于等于 20480, 直接读取整个响应体, 否则将使用 readChunks 进行分块读取, 优化内存占用和提升读取效率
	// 为了得到准确的 Content-Length , 在 newHTTPClient 中禁用了 gzip 压缩, 开启压缩后无法得到正确的 Content-Length
	const maxBodySizeForOneTimeRead = 20480
	if req.ContentLength >= 0 && req.ContentLength <= int64(maxBodySizeForOneTimeRead) {
		if _, err := io.Copy(&buf, req.Body); err != nil && err != io.EOF {
			// 如果 io.Copy 时发生错误, 暂不处理, 直接返回空 Body
			return nil
		}
	} else if err := readChunks(&buf, req.Body); err != nil {
		// 如果分块读取时发生错误, 暂不处理, 直接返回空 Body
		return nil
	}
	return buf.Bytes()
}

// readChunks 从 reader 中逐块读取数据并写入 writer
func readChunks(writer io.Writer, reader io.Reader) error {
	// 分块读取时, 每块为 4096 字节, 防止占用内存过大
	const chunkSize = 4096
	buffer := make([]byte, chunkSize)
	for n, err := reader.Read(buffer); n > 0 && err != io.EOF; n, err = reader.Read(buffer) {
		if _, err := writer.Write(buffer[:n]); err != nil {
			return err
		}
	}
	return nil
}
