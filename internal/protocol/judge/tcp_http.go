package judge

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/zhzyker/dismap/configs"
	"github.com/zhzyker/dismap/internal/proxy"
	"github.com/zhzyker/dismap/pkg/logger"
	"golang.org/x/text/encoding/simplifiedchinese"
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
)

func TcpHTTP(result map[string]interface{}, Args map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`^HTTP/\d.\d \d*`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "http"
		httpResult, httpErr := httpIdentifyResult(result, Args)
		if logger.DebugError(httpErr) {
			result["banner.string"] = "None"
			return true
		}
		result["banner.string"] = httpResult["http.title"].(string)
		u, err := url.Parse(httpResult["http.target"].(string))
		if err != nil {
			result["path"] = ""
		} else {
			result["path"] = u.Path
		}
		r := httpResult["http.result"].(string)
		c := fmt.Sprintf("[%s]", logger.Purple(httpResult["http.code"].(string)))
		if len(r) != 0 {
			result["identify.bool"] = true
			result["identify.string"] = fmt.Sprintf("%s %s", c, r)
			result["note"] = httpResult["http.target"].(string)
			return true
		} else {
			result["identify.bool"] = true
			result["identify.string"] = c
			result["note"] = httpResult["http.target"].(string)
			return true
		}
	}
	return false
}

func httpIdentifyResult(result map[string]interface{}, Args map[string]interface{}) (map[string]interface{}, error) {
	timeout := Args["FlagTimeout"].(int)
	var targetUrl string
	if Args["FlagUrl"].(string) != "" {
		targetUrl = Args["FlagUrl"].(string)
	} else{
		host := result["host"].(string)
		port := strconv.Itoa(result["port"].(int))
		add := net.JoinHostPort(host, port)
		if result["type"].(string) == "tcp" {
			if port == "80" {
				targetUrl = "http://" + host
			} else {
				targetUrl = "http://" + add
			}
		}
		if result["type"].(string) == "tls" {
			if port == "443" {
				targetUrl = "https://" + host
			} else {
				targetUrl = "https://" + add
			}
		}
	}

	var httpType string
	var httpCode string
	var httpResult string
	var httpUrl string
	var httpTitle string
	r, err := identify(targetUrl, timeout)
	if logger.DebugError(err) {
		return nil, err
	}
	for _, results := range r {
		httpType = results.Type
		httpCode = results.RespCode
		httpResult = results.Result
		httpUrl = results.Url
		httpTitle = results.Title
	}
	res := map[string]interface{}{
		"http.type": httpType,
		"http.code": httpCode,
		"http.result": httpResult,
		"http.target": httpUrl,
		"http.title": httpTitle,
	}
	return res, nil
}

type RespLab struct {
	Url string
	RespBody string
	RespHeader string
	RespStatusCode string
	RespTitle string
	faviconMd5 string
}

func getFaviconMd5(Url string, timeout int) string {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport {
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
	for key, value :=  range configs.DefaultHeader {
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

func defaultRequests(Url string, timeout int) ([]RespLab, error) {
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
	for key, value :=  range configs.DefaultHeader {
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
		for key, value :=  range configs.DefaultHeader {
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
			for key, value :=  range configs.DefaultHeader {
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
			RespData := []RespLab{
				{redirectUrl, responseBody, responseHeader, responseStatusCode, respTitle, faviconMd5},
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
		RespData := []RespLab{
			{redirectUrl, responseBody, responseHeader, responseStatusCode, respTitle, faviconMd5},
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
	RespData := []RespLab{
		{Url, responseBody, responseHeader, responseStatusCode, respTitle, faviconMd5},
	}
	return RespData, nil
}

func customRequests(Url string, timeout int, Method string, Path string, Header []string, Body string) ([]RespLab, error) {
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
		Transport: &http.Transport {
			TLSClientConfig:&tls.Config{InsecureSkipVerify: true},
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
	RespData := []RespLab{
		{Url, responseBody, responseHeader, responseStatusCode, respTitle, ""},
	}
	return RespData, nil
}

type IdentifyResult struct {
	Type     string
	RespCode string
	Result   string
	Url      string
	Title    string
}

func identify(url string, timeout int) ([]IdentifyResult, error) {
	var DefaultFavicon string
	var CustomFavicon string
	var DefaultTarget string
	var CustomTarget string
	var Favicon string
	var RequestRule string
	var RespTitle string
	var RespBody string
	var RespHeader string
	var RespCode string
	var DefaultRespTitle string
	var DefaultRespBody string
	var DefaultRespHeader string
	var DefaultRespCode string
	var CustomRespTitle string
	var CustomRespBody string
	var CustomRespHeader string
	var CustomRespCode string

	R, err := defaultRequests(url, timeout)
	if logger.DebugError(err) {
		return nil, err
	}
	for _, resp := range R {
		DefaultRespBody = resp.RespBody
		DefaultRespHeader = resp.RespHeader
		DefaultRespCode = resp.RespStatusCode
		DefaultRespTitle = resp.RespTitle
		DefaultTarget = resp.Url
		DefaultFavicon = resp.faviconMd5
	}
	// start identify
	var succes_type string
	var identify_result string
	type Identify_Result struct {
		Name string
		Rank int
		Type string
	}
	var IdentifyData []Identify_Result
	for _, rule := range configs.RuleData {
		if rule.Http.ReqMethod != "" {
			r, err := customRequests(url, timeout, rule.Http.ReqMethod, rule.Http.ReqPath, rule.Http.ReqHeader, rule.Http.ReqBody)
			if logger.DebugError(err) {
				return nil, err
			}

			for _, resp := range r {
				CustomRespBody = resp.RespBody
				CustomRespHeader = resp.RespHeader
				CustomRespCode = resp.RespStatusCode
				CustomRespTitle = resp.RespTitle
				CustomTarget = resp.Url
				CustomFavicon = resp.faviconMd5
			}
			url = CustomTarget
			Favicon = CustomFavicon
			RespBody = CustomRespBody
			RespHeader = CustomRespHeader
			RespCode = CustomRespCode
			RespTitle = CustomRespTitle
			// If the http request fails, then RespBody and RespHeader are both null
			// At this time, it is considered that the url does not exist
			if RespBody == RespHeader {
				continue
			}
			if rule.Mode == "" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "CustomRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "or" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "and" {
				index := 0
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						index = index + 1
					}
				}
				if index == 2 {
					IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
					RequestRule = "CustomRequest"
				}
			}
			if rule.Mode == "and|and" {
				index := 0
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						index = index + 1
					}
				}
				if index == 3 {
					IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
					RequestRule = "CustomRequest"
				}
			}
			if rule.Mode == "or|or" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "and|or" {
				grep := regexp.MustCompile("(.*)\\|(.*)\\|(.*)")
				all_type := grep.FindStringSubmatch(rule.Type)
				if len(regexp.MustCompile("header").FindAllStringIndex(all_type[1], -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == checkFavicon(Favicon, rule.Rule.InIcoMd5) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(all_type[1], -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == checkFavicon(Favicon, rule.Rule.InIcoMd5) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(all_type[1], -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "or|and" {
				grep := regexp.MustCompile("(.*)\\|(.*)\\|(.*)")
				all_type := grep.FindStringSubmatch(rule.Type)
				fmt.Println(all_type)
				if len(regexp.MustCompile("header").FindAllStringIndex(all_type[3], -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == checkFavicon(Favicon, rule.Rule.InIcoMd5) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(all_type[3], -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == checkFavicon(Favicon, rule.Rule.InIcoMd5) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(all_type[3], -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						succes_type = rule.Type
						continue
					}
				}
			}
		} else { // Default Request Result
			url = DefaultTarget
			Favicon = DefaultFavicon
			RespBody = DefaultRespBody
			RespHeader = DefaultRespHeader
			RespCode = DefaultRespCode
			RespTitle = DefaultRespTitle
			// If the http request fails, then RespBody and RespHeader are both null
			// At this time, it is considered that the url does not exist
			if RespBody == RespHeader {
				continue
			}
			if rule.Mode == "" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "or" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "and" {
				index := 0
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						index = index + 1
					}
				}
				if index == 2 {
					IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
					RequestRule = "DefaultRequest"
				}
			}
			if rule.Mode == "and|and" {
				index := 0
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						index = index + 1
					}
				}
				if index == 3 {
					IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
					RequestRule = "DefaultRequest"
				}
			}
			if rule.Mode == "or|or" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == true {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "and|or" {
				grep := regexp.MustCompile("(.*)\\|(.*)\\|(.*)")
				all_type := grep.FindStringSubmatch(rule.Type)
				fmt.Println(all_type)
				if len(regexp.MustCompile("header").FindAllStringIndex(all_type[1], -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == checkFavicon(Favicon, rule.Rule.InIcoMd5) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(all_type[1], -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == checkFavicon(Favicon, rule.Rule.InIcoMd5) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(all_type[1], -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "or|and" {
				grep := regexp.MustCompile("(.*)\\|(.*)\\|(.*)")
				all_type := grep.FindStringSubmatch(rule.Type)
				fmt.Println(all_type)
				if len(regexp.MustCompile("header").FindAllStringIndex(all_type[3], -1)) == 1 {
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == checkFavicon(Favicon, rule.Rule.InIcoMd5) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(all_type[3], -1)) == 1 {
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == checkFavicon(Favicon, rule.Rule.InIcoMd5) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(all_type[3], -1)) == 1 {
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == checkHeader(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if checkFavicon(Favicon, rule.Rule.InIcoMd5) == checkBody(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						IdentifyData = append(IdentifyData, Identify_Result {Name: rule.Name, Rank: rule.Rank, Type: rule.Type})
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
			}
		}
	}
	if RequestRule == "DefaultRequest" {
		RespBody = DefaultRespBody
		RespHeader = DefaultRespHeader
		RespCode = DefaultRespCode
		RespTitle = DefaultRespTitle
		url = DefaultTarget
	} else if RequestRule == "CustomRequest" {
		url = CustomTarget
		RespBody = CustomRespBody
		RespHeader = CustomRespHeader
		RespCode = CustomRespCode
		RespTitle = CustomRespTitle
	}
	for _, rs := range IdentifyData {
		switch rs.Rank {
		case 1:
			identify_result += "[" + logger.LightYellow(rs.Name) + "]"
		case 2:
			identify_result += "[" + logger.LightYellow(rs.Name) + "]"
		case 3:
			identify_result += "[" + logger.LightRed(rs.Name) + "]"
		}
	}
	r := strings.ReplaceAll(identify_result, "][", "] [")
	res := []IdentifyResult{{succes_type, RespCode, r,  url, RespTitle}}
	return res, nil
}

func checkHeader(url, responseHeader string, ruleHeader string, name string, title string, RespCode string) bool {
	grep := regexp.MustCompile("(?i)" + ruleHeader)
	if len(grep.FindStringSubmatch(responseHeader)) != 0 {
		//fmt.Print("[header] ")
		return true
	} else {
		return false
	}
}

func checkBody(url, responseBody string, ruleBody string, name string, title string, RespCode string) bool {
	grep := regexp.MustCompile("(?i)" + ruleBody)
	if len(grep.FindStringSubmatch(responseBody)) != 0 {
		//fmt.Print("[body] ")
		return true
	} else {
		return false
	}
}

func checkFavicon(Favicon, ruleFaviconMd5 string) bool {
	grep := regexp.MustCompile("(?i)" + ruleFaviconMd5)
	if len(grep.FindStringSubmatch(Favicon)) != 0 {
		return true
	} else {
		return false
	}
}