package http

import (
	"golang.org/x/net/html"
	"strings"
)

// getTitle 使用 golang.org/x/net/html 提取 title 为 string
func getTitle(htmlContent []byte) string {
	// 使用 html.Parse 解析 HTML 内容
	doc, _ := html.Parse(strings.NewReader(string(htmlContent)))
	var title string
	var findTitle func(*html.Node) // findTitle 是一个递归函数，用于遍历 HTML 结构查找 title
	findTitle = func(n *html.Node) {
		// 如果节点非空，节点类型为 ElementNode，节点标签为 "title"，并且有子节点
		if n != nil && n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil && title == ""; c = c.NextSibling {
			findTitle(c)
		}
	}
	// 调用 findTitle 函数查找 title
	findTitle(doc)
	return title
}
