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
	req, err := MakeDefaultRequest("https://www.baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	s, err := RequestSample(req, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IdentifyRules(s, time.Second))
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
	}, s, time.Second))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body",
		Mode: "or",
		Rule: rule.InStr{
			InHeader: "apache",
			InBody:   "<h1>test</h>",
			InIcoMd5: "",
		},
	}, s, time.Second))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body",
		Mode: "and",
		Rule: rule.InStr{
			InHeader: "server",
			InBody:   "<h1>test</h>",
			InIcoMd5: "",
		},
	}, s, time.Second))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body",
		Mode: "and",
		Rule: rule.InStr{
			InHeader: "[xxxxxx]",
			InBody:   "<h1>test</h>",
			InIcoMd5: "",
		},
	}, s, time.Second))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body|ico",
		Mode: "and|and",
		Rule: rule.InStr{
			InHeader: "server",
			InBody:   "<h1>test</h>",
			InIcoMd5: "md5123123",
		},
	}, s, time.Second))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body|ico",
		Mode: "or|and",
		Rule: rule.InStr{
			InHeader: "notfound",
			InBody:   "notfound",
			InIcoMd5: "md5123123",
		},
	}, s, time.Second))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body|ico",
		Mode: "and|or",
		Rule: rule.InStr{
			InHeader: "server",
			InBody:   "<h1>(xxxxx)</h>",
			InIcoMd5: "(xxxxxx)",
		},
	}, s, time.Second))

	t.Log(identifyRule(rule.RuleLab{
		Name: "test",
		Type: "header|body|ico",
		Mode: "or|or",
		Rule: rule.InStr{
			InHeader: "server",
			InBody:   "<h1>(xxxxx)</h>",
			InIcoMd5: "(xxxxxx)",
		},
	}, s, time.Second))
}

func Test_identifyRule(t *testing.T) {
	req, err := MakeDefaultRequest("http://xxxxx")
	if err != nil {
		t.Fatal(err)
	}
	s, err := RequestSample(req, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", s)
	r := rule.RuleLab{"BT.cn", "body|header|ico", "or", rule.InStr{"(请使用正确的入口登录面板|rm -f /www/server/panel/data/admin_path.pl|宝塔(.*)面板)", "(Set-Cookie: BT_COLL=)", "(9637ebd168435de51fea8193d2d89e39)"}, rule.ReqHttp{"", "", nil, ""}}
	t.Log(identifyRule(r, s, 5*time.Second))
}
