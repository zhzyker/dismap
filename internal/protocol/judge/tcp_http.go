package judge

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/zhzyker/dismap/configs"
	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/proxy"
	"github.com/zhzyker/dismap/pkg/logger"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func TcpHTTP(result *model.Result) bool {
	var buff = result.BannerB
	ok, err := regexp.Match(`^HTTP/\d.\d \d*`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result.Protocol = "http"
		httpResult, fpHints, httpErr := httpIdentifyResult(result)
		if logger.DebugError(httpErr) {
			result.Banner = "None"
			result.Identify = fpHints
			return true
		}
		result.Identify = fpHints
		result.Banner = httpResult.Title
		u, err := url.Parse(httpResult.Url)
		if err != nil {
			result.Path = ""
		} else {
			result.Path = u.Path
		}
		r := httpResult.Result
		c := fmt.Sprintf("[%s]", logger.Purple(httpResult.StatusCode))
		if len(r) != 0 {
			result.IdentifyBool = true
			result.IdentifyStr = fmt.Sprintf("%s %s", c, r)
			result.Note = httpResult.Url
			return true
		} else {
			result.IdentifyBool = true
			result.IdentifyStr = c
			result.Note = httpResult.Url
			return true
		}
	}
	return false
}

func httpIdentifyResult(result *model.Result) (*model.HttpResult, []model.HintFinger, error) {
	timeout := flag.Timeout
	var targetUrl string
	if flag.InUrl != "" {
		targetUrl = flag.InUrl
	} else {
		host := result.Host
		port := strconv.Itoa(result.Port)
		add := net.JoinHostPort(host, port)
		if result.Type == "tcp" {
			if port == "80" {
				targetUrl = "http://" + host
			} else {
				targetUrl = "http://" + add
			}
		}
		if result.Type == "tls" {
			if port == "443" {
				targetUrl = "https://" + host
			} else {
				targetUrl = "https://" + add
			}
		}
	}

	r, hint, err := identify(targetUrl, timeout)
	if logger.DebugError(err) {
		return nil, nil, err
	}
	return r, hint, nil
}

var RuleFuncs = map[string]func(bool, bool) bool{
	"and": checkRuleAnd,
	"or":  checkRuleOr,
}
var CheckFuncs = map[string]func(*model.HttpResult, *configs.RuleLab) bool{
	"body":   checkBody,
	"header": checkHeader,
	"ico":    checkFavicon,
}

func getFaviconMd5(Url string, timeout int) string {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	Url = Url + "/favicon.ico"
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return ""
	}
	for key, value := range configs.DefaultHeader {
		req.Header.Set(key, value)
	}
	//req.Header.Set("Accept-Language", "zh,zh-TW;q=0.9,en-US;q=0.8,en;q=0.7,zh-CN;q=0.6")
	//req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	//req.Header.Set("Cookie", "rememberMe=int")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	hash := md5.Sum(bodyBytes)
	md5 := fmt.Sprintf("%x", hash)
	return md5
}

func defaultRequests(Url string, timeout int) (*model.HttpResult, error) {
	var redirectUrl string
	var respTitle string
	var responseHeader string
	var responseBody string
	var responseStatusCode string
	var res []string

	req, err := http.NewRequest("GET", Url, nil)
	if logger.DebugError(err) {
		return nil, err
	}
	// set requests header
	for key, value := range configs.DefaultHeader {
		req.Header.Set(key, value)
	}
	resp, err := proxy.ConnProxyHttp(req, timeout)
	if logger.DebugError(err) {
		return nil, err
	}
	defer resp.Body.Close()
	// get response status code
	var statusCode = resp.StatusCode
	responseStatusCode = strconv.Itoa(statusCode)
	// -------------------------------------------------------------------------------
	// When the http request is 302 or other 30x,
	// Need to intercept the request and get the return status code for display,
	// Send the request again according to the redirect url
	// In the custom request, the return status code is not checked
	// -------------------------------------------------------------------------------
	if len(regexp.MustCompile("30").FindAllStringIndex(responseStatusCode, -1)) == 1 {
		redirectPath := resp.Header.Get("Location")
		if len(regexp.MustCompile("http").FindAllStringIndex(redirectPath, -1)) == 1 {
			redirectUrl = redirectPath
		} else {
			if Url[len(Url)-1:] == "/" {
				redirectUrl = Url + redirectPath
			}
			redirectUrl = Url + "/" + redirectPath
		}
		req, err := http.NewRequest("GET", redirectUrl, nil)
		if err != nil {
			return nil, err
		}
		for key, value := range configs.DefaultHeader {
			req.Header.Set(key, value)
		}
		resp, err := proxy.ConnProxyHttp(req, timeout)
		if logger.DebugError(err) {
			return nil, err
		}
		defer resp.Body.Close()
		// Solve the problem of two 30x jumps
		var twoStatusCode = resp.StatusCode
		responseStatusCodeTwo := strconv.Itoa(twoStatusCode)
		if len(regexp.MustCompile("30").FindAllStringIndex(responseStatusCodeTwo, -1)) == 1 {
			redirectPath := resp.Header.Get("Location")
			if len(regexp.MustCompile("http").FindAllStringIndex(redirectPath, -1)) == 1 {
				redirectUrl = redirectPath
			} else {
				redirectUrl = Url + redirectPath
			}
			req, err := http.NewRequest("GET", redirectUrl, nil)
			if err != nil {
				return nil, err
			}
			for key, value := range configs.DefaultHeader {
				req.Header.Set(key, value)
			}
			resp, err := proxy.ConnProxyHttp(req, timeout)
			if logger.DebugError(err) {
				return nil, err
			}
			defer resp.Body.Close()
			// get response body for string
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			responseBody = string(bodyBytes)
			// Solve the problem of garbled body codes with unmatched numbers
			if !utf8.Valid(bodyBytes) {
				data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(bodyBytes)
				responseBody = string(data)
			}
			// Get Response title
			grepTitle := regexp.MustCompile("<title>(.*)</title>")
			if len(grepTitle.FindStringSubmatch(responseBody)) != 0 {
				respTitle = grepTitle.FindStringSubmatch(responseBody)[1]
			} else {
				respTitle = "None"
			}
			// get response header for string
			for name, values := range resp.Header {
				for _, value := range values {
					res = append(res, fmt.Sprintf("%s: %s", name, value))
				}
			}
			for _, re := range res {
				responseHeader += re + "\n"
			}
			faviconMd5 := getFaviconMd5(Url, timeout)
			RespData := &model.HttpResult{
				Url:        redirectUrl,
				Body:       responseBody,
				Header:     responseHeader,
				StatusCode: responseStatusCode,
				Title:      respTitle,
				FaviconMd5: faviconMd5,
			}
			return RespData, nil
		}
		// get response body for string
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		responseBody = string(bodyBytes)
		// Solve the problem of garbled body codes with unmatched numbers
		if !utf8.Valid(bodyBytes) {
			data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(bodyBytes)
			responseBody = string(data)
		}
		// Get Response title
		grepTitle := regexp.MustCompile("<title>(.*)</title>")
		if len(grepTitle.FindStringSubmatch(responseBody)) != 0 {
			respTitle = grepTitle.FindStringSubmatch(responseBody)[1]
		} else {
			respTitle = "None"
		}
		// get response header for string
		for name, values := range resp.Header {
			for _, value := range values {
				res = append(res, fmt.Sprintf("%s: %s", name, value))
			}
		}
		for _, re := range res {
			responseHeader += re + "\n"
		}
		faviconMd5 := getFaviconMd5(Url, timeout)
		RespData := &model.HttpResult{
			Url:        redirectUrl,
			Body:       responseBody,
			Header:     responseHeader,
			StatusCode: responseStatusCode,
			Title:      respTitle,
			FaviconMd5: faviconMd5,
		}

		return RespData, nil
	}
	// get response body for string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	responseBody = string(bodyBytes)
	// Solve the problem of garbled body codes with unmatched numbers
	if !utf8.Valid(bodyBytes) {
		data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(bodyBytes)
		responseBody = string(data)
	}
	// Get Response title
	grepTitle := regexp.MustCompile("<title>(.*)</title>")
	if len(grepTitle.FindStringSubmatch(responseBody)) != 0 {
		respTitle = grepTitle.FindStringSubmatch(responseBody)[1]
	} else {
		respTitle = "None"
	}
	// get response header for string
	for name, values := range resp.Header {
		for _, value := range values {
			res = append(res, fmt.Sprintf("%s: %s", name, value))
		}
	}
	for _, re := range res {
		responseHeader += re + "\n"
	}
	faviconMd5 := getFaviconMd5(Url, timeout)
	RespData := &model.HttpResult{
		Url:        Url,
		Body:       responseBody,
		Header:     responseHeader,
		StatusCode: responseStatusCode,
		Title:      respTitle,
		FaviconMd5: faviconMd5,
	}
	return RespData, nil
}

func customRequests(Url string, timeout int, Method string, Path string, Header []string, Body string) (*model.HttpResult, error) {
	var respTitle string
	// Splicing Custom Path
	u, err := url.Parse(Url)
	u.Path = path.Join(u.Path, Path)
	Url = u.String()
	if strings.HasSuffix(Path, "/") {
		Url = Url + "/"
	}
	// Send Http requests
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	bodyByte := bytes.NewBuffer([]byte(Body))
	req, err := http.NewRequest(Method, Url, bodyByte)
	if logger.DebugError(err) {
		return nil, err
	}

	// Set Requests Headers
	for _, header := range Header {
		grepKey := regexp.MustCompile("(.*): ")
		var headerKey = grepKey.FindStringSubmatch(header)[1]
		grepValue := regexp.MustCompile(": (.*)")
		var headerValue = grepValue.FindStringSubmatch(header)[1]
		req.Header.Set(headerKey, headerValue)
	}
	resp, err := client.Do(req)
	if logger.DebugError(err) {
		return nil, err
	}
	defer resp.Body.Close()
	// Get Response Body for string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var responseBody = string(bodyBytes)
	// Solve the problem of garbled body codes with unmatched numbers
	if !utf8.Valid(bodyBytes) {
		data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(bodyBytes)
		responseBody = string(data)
	}
	// Get Response title
	grepTitle := regexp.MustCompile("<title>(.*)</title>")
	if len(grepTitle.FindStringSubmatch(responseBody)) != 0 {
		respTitle = grepTitle.FindStringSubmatch(responseBody)[1]
	} else {
		respTitle = "None"
	}
	// Get Response Header for string
	var res []string
	for name, values := range resp.Header {
		for _, value := range values {
			res = append(res, fmt.Sprintf("%s: %s", name, value))
		}
	}
	var responseHeader string
	for _, re := range res {
		responseHeader += re + "\n"
	}
	// get response status code
	var statusCode = resp.StatusCode
	responseStatusCode := strconv.Itoa(statusCode)
	faviconMd5 := getFaviconMd5(Url, timeout)
	RespData := &model.HttpResult{
		Url:        Url,
		Body:       responseBody,
		Header:     responseHeader,
		StatusCode: responseStatusCode,
		Title:      respTitle,
		FaviconMd5: faviconMd5,
	}
	return RespData, nil
}

func identify(url string, timeout int) (*model.HttpResult, []model.HintFinger, error) {
	var checkResp *model.HttpResult

	defaultResp, err := defaultRequests(url, timeout)
	if err != nil {
		return nil, nil, err
	}

	var Hints []model.HintFinger

	for _, rule := range configs.RuleData {
		// 如果规则需要自定义请求，那么就请求一下
		if rule.Http.ReqMethod != "" {
			r, err := customRequests(url, timeout, rule.Http.ReqMethod, rule.Http.ReqPath, rule.Http.ReqHeader, rule.Http.ReqBody)
			if err != nil {
				return nil, nil, err
			}
			checkResp = r
		} else {
			// 否则使用默认数据
			// Default Request Result
			checkResp = defaultResp
		}
		// If the http request fails, then RespBody and RespHeader are both null
		// At this time, it is considered that the url does not exist
		if checkResp.Body == checkResp.Header {
			continue
		}
		// 开始判断
		modes := strings.Split(rule.Mode, "|")
		strs := strings.Split(rule.Type, "|")
		if modes[0] == "" || len(strs) == 1 {
			if ff, ok := CheckFuncs[rule.Type]; ok {
				if ff(checkResp, rule) {
					Hints = append(Hints, model.HintFinger{Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
					continue
				}
			}
		} else if len(modes) == 1 {
			if RuleFuncs[modes[0]](CheckFuncs[strs[0]](checkResp, rule), CheckFuncs[strs[1]](checkResp, rule)) {
				Hints = append(Hints, model.HintFinger{Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
				continue
			}
		} else if len(modes) == 2 {
			status := CheckFuncs[strs[0]](checkResp, rule)
			for index, _rt := range modes {
				status = RuleFuncs[_rt](status, CheckFuncs[strs[index+1]](checkResp, rule))
				if status {
					Hints = append(Hints, model.HintFinger{Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
					break
				}
			}
			continue
		}
	}

	for index := range Hints {
		Hints[index].Source = "dismap"
	}

	for _, hint := range Hints {
		switch hint.Rank {
		case 1:
			checkResp.Result += "[" + logger.LightYellow(hint.Name) + "]"
		case 2:
			checkResp.Result += "[" + logger.LightYellow(hint.Name) + "]"
		case 3:
			checkResp.Result += "[" + logger.LightRed(hint.Name) + "]"
		}
	}

	checkResp.Result = strings.ReplaceAll(checkResp.Result, "][", "] [")
	return checkResp, Hints, nil
}

func checkRuleAnd(status1, status2 bool) bool {
	return status1 && status2
}

func checkRuleOr(status1, status2 bool) bool {
	return status1 || status2
}

func checkHeader(resp *model.HttpResult, rule *configs.RuleLab) bool {
	grep := regexp.MustCompile("(?i)" + rule.Rule.InHeader)
	if len(grep.FindStringSubmatch(resp.Header)) != 0 {
		//fmt.Print("[header] ")
		return true
	} else {
		return false
	}
}

func checkBody(resp *model.HttpResult, rule *configs.RuleLab) bool {
	grep := regexp.MustCompile("(?i)" + rule.Rule.InBody)
	if len(grep.FindStringSubmatch(resp.Body)) != 0 {
		//fmt.Print("[body] ")
		return true
	} else {
		return false
	}
}

func checkFavicon(resp *model.HttpResult, rule *configs.RuleLab) bool {
	grep := regexp.MustCompile("(?i)" + rule.Rule.InIcoMd5)
	if len(grep.FindStringSubmatch(resp.FaviconMd5)) != 0 {
		return true
	} else {
		return false
	}
}
