package http

import (
	"regexp"
)

// getTitle 从 body 中正则提取 title
func getTitle(body string) string {
	match := regexp.MustCompile(`(?i)<title>\s*(.+?)\s*</title>`).FindStringSubmatch(body)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
