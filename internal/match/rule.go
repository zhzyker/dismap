package match

import (
	"encoding/json"
	"github.com/zhzyker/dismap/pkg/logger"
	"io/ioutil"
)

// Rule 规则结构体，包含规则的各种信息和规则匹配
type Rule struct {
	RuleID   string        `json:"rule_id"`
	Type     string        `json:"type"`
	Name     string        `json:"name"`
	Level    string        `json:"level"`
	Protocol string        `json:"protocol"`
	Rules    [][]RuleMatch `json:"rules"`
}

// RuleMatch 规则匹配，包含规则匹配的条件和内容
type RuleMatch struct {
	Match   string `json:"match"`
	Content string `json:"content"`
}

var rulesCache []Rule // ReadRules 读取规则文件返回规则

func ReadRules() []Rule {
	if len(rulesCache) > 0 {
		return rulesCache
	}
	// 读取 JSON 文件
	jsonData, err := ioutil.ReadFile("config/rules.json")
	if err != nil {
		logger.ERR("Failed to read rules.json file: " + err.Error())
		return nil
	}
	// 解码 JSON 数据到 rules
	var rules []Rule
	if err := json.Unmarshal(jsonData, &rules); err != nil {
		logger.ERR("Failed to unmarshal rules.json file: " + err.Error())
		return nil
	}
	rulesCache = rules
	return rules
}
