package dismap

import (
	"regexp"
	"strings"
	"time"

	"github.com/zhzyker/dismap/pkg/slice"
	"github.com/zhzyker/dismap/rule"
)

func IdentifyRules(sample *Sample, timeout time.Duration) []string {
	resutls := make([]string, 0)
	for _, rule := range rule.RuleData {
		if identifyRule(rule, sample, timeout) {
			resutls = append(resutls, rule.Name)
		}
	}
	return resutls
}

func identifyRule(rule rule.RuleLab, sample *Sample, timeout time.Duration) bool {
	// make custom sample
	if rule.Http.ReqMethod != "" {
		req, err := MakeCustomRequest(sample.Url, rule.Http.ReqMethod, rule.Http.ReqPath, rule.Http.ReqHeader, rule.Http.ReqBody)
		if err == nil {
			if s, err := RequestSample(req, timeout); err == nil {
				sample = s
			}
		}
	}
	types := slice.RemoveDuplicationSort(strings.Split(rule.Type, "|"))
	operators := make([]string, 0)
	if rule.Mode != "" {
		operators = strings.Split(rule.Mode, "|")
	}

	diff := len(types) - len(operators)
	if diff != 1 {
		if len(operators) != 1 {
			return false
		}
		// 支持操作符仅为一个时
		for i := 0; i < diff; i++ {
			operators = append(operators, operators[0])
		}
	}

	var res bool
	for i := range types {
		if i == 0 {
			res = checkRuleType(types[i], rule, sample)
			continue
		}
		switch operators[i-1] {
		case "or":
			if res {
				continue
			} else {
				res = res || checkRuleType(types[i], rule, sample)
			}
		case "and":
			if !res {
				continue
			} else {
				res = res && checkRuleType(types[i], rule, sample)
			}
		default:
			res = false
		}
	}
	return res
}

func checkRuleType(key string, rule rule.RuleLab, sample *Sample) bool {
	switch key {
	case "header":
		return checkInContent(rule.Rule.InHeader, sample.Header)
	case "body":
		return checkInContent(rule.Rule.InBody, sample.Body)
	case "ico":
		return checkInContent(rule.Rule.InIcoMd5, sample.FaviconMd5)
	}
	return false
}

func checkInContent(reg, content string) bool {
	grep := regexp.MustCompile("(?i)" + reg)
	if len(grep.FindStringSubmatch(content)) != 0 {
		return true
	} else {
		return false
	}
}
