package match

import (
	"fmt"
	"github.com/zhzyker/dismap/pkg/logger"
	"github.com/zhzyker/dismap/pkg/requet/http"
	"regexp"
)

// switchMatch 用于判断响应是否匹配规则
func switchMatch(match string, content string, res http.HttpResult) bool {
	switch match {
	case "body":
		return regexp.MustCompile(content).MatchString(res.Body)
	case "header":
		return regexp.MustCompile(content).MatchString(res.Headers)
	case "favicon":
		//return regexp.MustCompile(content).MatchString(res.Favicon)
	}
	return false
}

// IdentifyResource 通过加载指纹库来识别 URL 对应的指纹结果
func IdentifyResource(url string) ([]byte, error) {
	logger.INF("Load rules to start identifying target: " + url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var matches []addMatch // 用于存储所有匹配的指纹
	// 从文件中读取规则并迭代匹配
	for _, r := range ReadRules() {
		for _, content := range r.Rules {
			matched := true
			for _, c := range content {
				matched = matched && switchMatch(c.Match, c.Content, res)
			}
			// 如果指纹匹配成功，则将其添加到匹配结果中
			if matched {
				matches = append(matches, addMatch{I: r.Name})
				logger.DBG(fmt.Sprintf("Hit to name: %s, Rule id: %s The rule is: %s", r.Name, r.RuleID, content))
			}
		}
	}
	return identifyResult(matches, res, url), nil
}
