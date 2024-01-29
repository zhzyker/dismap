package http

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
)

func getFavicon(url string, client *http.Client, isFavicon bool) []byte {
	if !isFavicon {
		return nil
	}
	req, err := client.Get(url + "/favicon.ico")
	if err != nil || req.StatusCode != http.StatusOK {
		return nil // 处理 HTTP 请求发送失败或状态码不是 200 OK 的情况
	}
	defer func() {
		_ = req.Body.Close() // 处理关闭 Body 过程中的错误，可以记录日志等
	}()
	hash := md5.New()
	_, _ = io.Copy(hash, req.Body)
	return []byte(hex.EncodeToString(hash.Sum(nil)))
}
