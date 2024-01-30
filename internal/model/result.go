package model

import "time"

type Result struct {
	Date         time.Time
	Status       string
	Banner       string
	BannerB      []byte
	Protocol     string // HTTP / MYSQL ...
	Type         string // TCP / UDP
	Host         string
	Port         int
	Uri          string
	Note         string
	Path         string
	IdentifyBool bool
	IdentifyStr  string
	Identify     []HintFinger
	HttpResult   *HttpResult
	Hint         bool
}

type HttpResult struct {
	Url            string              // URL
	Body           string              // Body
	Server         string              // WebServer
	Header         string              // Header
	HeaderM        map[string][]string // Header map
	StatusCode     string              // 状态码
	Title          string              // 标题
	FaviconMd5     string              // 图标Md5
	FaviconMM3Hash string              // mm3Hash (指纹专用)
	Type           string              // HTTP/HTTPs
	Result         string              // dismap 原版的信息存储位置
}

type HintFinger struct {
	Name        string
	Version     string
	Description string
	Rank        int
	Type        string
	Source      string
}
