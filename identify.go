package dismap

import (
	"regexp"
	"strings"

	"github.com/zhzyker/dismap/pkg/slice"
	"github.com/zhzyker/dismap/rule"
)

func IdentifyRules(sample *Sample) []string {
	resutls := make([]string, 0)
	for _, rule := range rule.RuleData {
		if identifyRule(rule, sample) {
			resutls = append(resutls, rule.Name)
		}
	}
	return resutls
}

func identifyRule(rule rule.RuleLab, sample *Sample) bool {
	// TODO
	// custom make sample
	if rule.Http.ReqMethod != "" {
	}
	types := slice.RemoveDuplicationSort(strings.Split(rule.Type, "|"))
	operators := make([]string, 0)
	if rule.Mode != "" {
		operators = strings.Split(rule.Mode, "|")
	}
	if len(types)-len(operators) != 1 {
		return false
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
