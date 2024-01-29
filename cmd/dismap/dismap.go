package main

import (
	"fmt"
	"github.com/zhzyker/dismap/internal/match"
	"github.com/zhzyker/dismap/pkg/requet/http"
	"net/url"
)

/*
_____________
______  /__(_)_____________ _________ ________
_  __  /__  /__  ___/_  __ `__ \  __ `/__  __ \
/ /_/ / _  / _(__  )_  / / / / / /_/ /__  /_/ /
\__,_/  /_/  /____/ /_/ /_/ /_/\__,_/ _  .___/
                                        /_/
  author: zhzyker && Nemophllist
  from: https://github.com/zhzyker/dismap
*/

func main() {
	url := "https://www.baidu.com"
	response, _ := http.Request(url, 10, true, MustParseProxyURL("http://127.0.0.1:8080"))
	var res []byte
	res, _ = match.IdentifyResource(response)
	fmt.Println(string(res))
}

// MustParseProxyURL 用于解析代理地址，发生错误时 panic
func MustParseProxyURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return parsedURL
}
