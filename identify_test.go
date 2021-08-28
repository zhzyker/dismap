package dismap

import (
	"testing"
	"time"

	"github.com/zhzyker/dismap/rule"
)

func Test_checkInContent(t *testing.T) {
	t.Log(checkInContent("(7dbe9acc2ab6e64d59fa67637b1239df|asdasdasd)", "7dbe9acc2ab6e64d59fa67637b1239df"))

}

func Test_IdentifyRules(t *testing.T) {
	s, err := RequestSample("https://www.baidu.com", 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IdentifyRules(s))
}

func Test_identifyRules(t *testing.T) {
	s := &Sample{
		Header:     "server",
		Body:       "<h1>test</h>",
		FaviconMd5: "md5123123",
	}
	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header",
		Mode: "",
		Rule: rule.InStr{
			InHeader: "server",
			InBody:   "",
			InIcoMd5: "",
		},
	}, s))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body",
		Mode: "or",
		Rule: rule.InStr{
			InHeader: "apache",
			InBody:   "<h1>test</h>",
			InIcoMd5: "",
		},
	}, s))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body",
		Mode: "and",
		Rule: rule.InStr{
			InHeader: "server",
			InBody:   "<h1>test</h>",
			InIcoMd5: "",
		},
	}, s))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body",
		Mode: "and",
		Rule: rule.InStr{
			InHeader: "[xxxxxx]",
			InBody:   "<h1>test</h>",
			InIcoMd5: "",
		},
	}, s))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body|ico",
		Mode: "and|and",
		Rule: rule.InStr{
			InHeader: "server",
			InBody:   "<h1>test</h>",
			InIcoMd5: "md5123123",
		},
	}, s))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body|ico",
		Mode: "or|and",
		Rule: rule.InStr{
			InHeader: "notfound",
			InBody:   "notfound",
			InIcoMd5: "md5123123",
		},
	}, s))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body|ico",
		Mode: "and|or",
		Rule: rule.InStr{
			InHeader: "server",
			InBody:   "<h1>(xxxxx)</h>",
			InIcoMd5: "(xxxxxx)",
		},
	}, s))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body|ico",
		Mode: "or|or",
		Rule: rule.InStr{
			InHeader: "server",
			InBody:   "<h1>(xxxxx)</h>",
			InIcoMd5: "(xxxxxx)",
		},
	}, s))
}
