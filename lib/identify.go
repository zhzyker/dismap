package lib

import (
	"../config"
	"fmt"
	"regexp"
)

type IdentifyResult struct {
	Type string
	RespCode string
	Result string
	ResultNc string
	Url string
	Title string
}

func Identify(url string, timeout int) []IdentifyResult {
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
	for _, resp := range DefaultRequests(url, timeout) { // Default Request
		DefaultRespBody = resp.RespBody
		DefaultRespHeader = resp.RespHeader
		DefaultRespCode = resp.RespStatusCode
		DefaultRespTitle = resp.RespTitle
		DefaultTarget = resp.Url
		DefaultFavicon = resp.FaviconMd5
	}
	// start identify
	var identify_data []string
	var succes_type string
	for _, rule := range config.RuleData {
		if rule.Http.ReqMethod != "" { // Custom Request Result
			for _, resp := range CustomRequests(url, timeout, rule.Http.ReqMethod, rule.Http.ReqPath, rule.Http.ReqHeader, rule.Http.ReqBody) {
				CustomRespBody = resp.RespBody
				CustomRespHeader = resp.RespHeader
				CustomRespCode = resp.RespStatusCode
				CustomRespTitle = resp.RespTitle
				CustomTarget = resp.Url
				CustomFavicon = resp.FaviconMd5
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
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "CustomRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "or" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "and" {
				index := 0
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						index = index + 1
					}
				}
				if index == 2 {
					identify_data = append(identify_data, rule.Name)
					RequestRule = "CustomRequest"
				}
			}
			if rule.Mode == "and|and" {
				index := 0
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						index = index + 1
					}
				}
				if index == 3 {
					identify_data = append(identify_data, rule.Name)
					RequestRule = "CustomRequest"
				}
			}
			if rule.Mode == "or|or" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						identify_data = append(identify_data, rule.Name)
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
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == check_favicon(Favicon, rule.Rule.InIcoMd5) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(all_type[1], -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == check_favicon(Favicon, rule.Rule.InIcoMd5) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(all_type[1], -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
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
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == check_favicon(Favicon, rule.Rule.InIcoMd5) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(all_type[3], -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == check_favicon(Favicon, rule.Rule.InIcoMd5) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(all_type[3], -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						succes_type = rule.Type
						continue
					}
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
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
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "or" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
			}
			if rule.Mode == "and" {
				index := 0
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						index = index + 1
					}
				}
				if index == 2 {
					identify_data = append(identify_data, rule.Name)
					RequestRule = "DefaultRequest"
				}
			}
			if rule.Mode == "and|and" {
				index := 0
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						index = index + 1
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						index = index + 1
					}
				}
				if index == 3 {
					identify_data = append(identify_data, rule.Name)
					RequestRule = "DefaultRequest"
				}
			}
			if rule.Mode == "or|or" {
				if len(regexp.MustCompile("header").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == true {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(rule.Type, -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == true {
						identify_data = append(identify_data, rule.Name)
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
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == check_favicon(Favicon, rule.Rule.InIcoMd5) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(all_type[1], -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == check_favicon(Favicon, rule.Rule.InIcoMd5) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(all_type[1], -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
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
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) == check_favicon(Favicon, rule.Rule.InIcoMd5) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("body").FindAllStringIndex(all_type[3], -1)) == 1 {
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) == check_favicon(Favicon, rule.Rule.InIcoMd5) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
				if len(regexp.MustCompile("ico").FindAllStringIndex(all_type[3], -1)) == 1 {
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == check_header(url, RespHeader, rule.Rule.InHeader, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
					if check_favicon(Favicon, rule.Rule.InIcoMd5) == check_body(url, RespBody, rule.Rule.InBody, rule.Name, RespTitle, RespCode) {
						identify_data = append(identify_data, rule.Name)
						RequestRule = "DefaultRequest"
						succes_type = rule.Type
						continue
					}
				}
			}
		}
	}
	// identify
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
	var identify_result string
	var identify_result_nocolor string
	for _, result := range identify_data {
		if sysarch == "windows" {
			identify_result += "["+result+"]"+" "
		} else {
			identify_result += "["+Yellow(result)+"]"+" "
		}
	}
	for _, result := range identify_data {
		identify_result_nocolor += "["+result+"]"+" "
	}

	Result := []IdentifyResult{
		{succes_type, RespCode, identify_result, identify_result_nocolor, url, RespTitle},
	}
	return Result
}

func check_header(url, response_header string, rule_header string, name string, title string, RespCode string) bool {
	grep := regexp.MustCompile("(?i)"+rule_header)
	if len(grep.FindStringSubmatch(response_header)) != 0 {
		//fmt.Print("[header] ")
		return true
	} else {
		return false
	}
}

func check_body(url, response_body string, rule_body string, name string, title string, RespCode string) bool {
	grep := regexp.MustCompile("(?i)"+rule_body)
	if len(grep.FindStringSubmatch(response_body)) != 0 {
		//fmt.Print("[body] ")
		return true
	} else {
		return false
	}
}

func check_favicon(Favicon, rule_favicon_md5 string) bool {
	grep := regexp.MustCompile("(?i)"+rule_favicon_md5)
	if len(grep.FindStringSubmatch(Favicon)) != 0 {
		// fmt.Print("url")
		return true
	} else {
		return false
	}
}
