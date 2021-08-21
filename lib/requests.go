package lib

import (
	"github.com/zhzyker/dismap/config"
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type RespLab struct {
	Url string
	RespBody string
	RespHeader string
	RespStatusCode string
	RespTitle string
	FaviconMd5 string
}

func FaviconMd5(Url string, timeout int, Path string) string {
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
	for key, value :=  range config.DefaultHeader {
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
	body_bytes, err := ioutil.ReadAll(resp.Body)
	hash := md5.Sum(body_bytes)
	md5 := fmt.Sprintf("%x", hash)
	return md5
}

func DefaultRequests(Url string, timeout int) []RespLab {
	var redirect_url string
	var resp_title string
	var response_header string
	var response_body string
	var response_status_code string
	var res []string
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport {
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil
	}
	// set requests header
	for key, value :=  range config.DefaultHeader {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	// get response status code
	var status_code = resp.StatusCode
	response_status_code = strconv.Itoa(status_code)
	// -------------------------------------------------------------------------------
	// When the http request is 302 or other 30x,
	// Need to intercept the request and get the return status code for display,
	// Send the request again according to the redirect url
	// In the custom request, the return status code is not checked
	// -------------------------------------------------------------------------------
	if len(regexp.MustCompile("30").FindAllStringIndex(response_status_code, -1)) == 1 {
		redirect_path := resp.Header.Get("Location")
		if len(regexp.MustCompile("http").FindAllStringIndex(redirect_path, -1)) == 1 {
			redirect_url = redirect_path
		} else {
			redirect_url = Url + redirect_path
		}
		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport {
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req, err := http.NewRequest("GET", redirect_url, nil)
		if err != nil {
			return nil
		}
		for key, value :=  range config.DefaultHeader {
			req.Header.Set(key, value)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil
		}
		defer resp.Body.Close()
		// Solve the problem of two 30x jumps
		var two_status_code = resp.StatusCode
		response_status_code_two := strconv.Itoa(two_status_code)
		if len(regexp.MustCompile("30").FindAllStringIndex(response_status_code_two, -1)) == 1 {
			redirect_path := resp.Header.Get("Location")
			if len(regexp.MustCompile("http").FindAllStringIndex(redirect_path, -1)) == 1 {
				redirect_url = redirect_path
			} else {
				redirect_url = Url + redirect_path
			}
			client := &http.Client{
				Timeout: time.Duration(timeout) * time.Second,
				Transport: &http.Transport {
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			req, err := http.NewRequest("GET", redirect_url, nil)
			if err != nil {
				return nil
			}
			for key, value :=  range config.DefaultHeader {
				req.Header.Set(key, value)
			}
			resp, err := client.Do(req)
			if err != nil {
				return nil
			}
			defer resp.Body.Close()
			// get response body for string
			body_bytes, err := ioutil.ReadAll(resp.Body)
			response_body = string(body_bytes)
			// Solve the problem of garbled body codes with unmatched numbers
			if !utf8.Valid(body_bytes) {
				data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(body_bytes)
				response_body = string(data)
			}
			// Get Response title
			grep_title := regexp.MustCompile("<title>(.*)</title>")
			if len(grep_title.FindStringSubmatch(response_body)) != 0 {
				resp_title = grep_title.FindStringSubmatch(response_body)[1]
			} else {
				resp_title = "None"
			}
			// get response header for string
			for name, values := range resp.Header {
				for _, value := range values {
					res = append(res, fmt.Sprintf("%s: %s", name, value))
				}
			}
			for _, re := range res {
				response_header += re + "\n"
			}
			faviconmd5 := FaviconMd5(Url, timeout, "")
			RespData := []RespLab{
				{redirect_url, response_body, response_header, response_status_code, resp_title, faviconmd5},
			}
			return RespData
		}
		// get response body for string
		body_bytes, err := ioutil.ReadAll(resp.Body)
		response_body = string(body_bytes)
		// Solve the problem of garbled body codes with unmatched numbers
		if !utf8.Valid(body_bytes) {
			data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(body_bytes)
			response_body = string(data)
		}
		// Get Response title
		grep_title := regexp.MustCompile("<title>(.*)</title>")
		if len(grep_title.FindStringSubmatch(response_body)) != 0 {
			resp_title = grep_title.FindStringSubmatch(response_body)[1]
		} else {
			resp_title = "None"
		}
		// get response header for string
		for name, values := range resp.Header {
			for _, value := range values {
				res = append(res, fmt.Sprintf("%s: %s", name, value))
			}
		}
		for _, re := range res {
			response_header += re + "\n"
		}
		faviconmd5 := FaviconMd5(Url, timeout, "")
		RespData := []RespLab{
			{redirect_url, response_body, response_header, response_status_code, resp_title, faviconmd5},
		}
		return RespData
	}

	// get response body for string
	body_bytes, err := ioutil.ReadAll(resp.Body)
	response_body = string(body_bytes)
	// Solve the problem of garbled body codes with unmatched numbers
	if !utf8.Valid(body_bytes) {
		data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(body_bytes)
		response_body = string(data)
	}

	// Get Response title
	grep_title := regexp.MustCompile("<title>(.*)</title>")
	if len(grep_title.FindStringSubmatch(response_body)) != 0 {
		resp_title = grep_title.FindStringSubmatch(response_body)[1]
	} else {
		resp_title = "None"
	}
	// get response header for string
	for name, values := range resp.Header {
		for _, value := range values {
			res = append(res, fmt.Sprintf("%s: %s", name, value))
		}
	}
	for _, re := range res {
		response_header += re + "\n"
	}
	faviconmd5 := FaviconMd5(Url, timeout, "")
	RespData := []RespLab{
		{Url, response_body, response_header, response_status_code, resp_title, faviconmd5},
	}
	return RespData
}

func CustomRequests(Url string, timeout int, Method string, Path string, Header []string, Body string) []RespLab {
	var resp_title string
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
	body_byte := bytes.NewBuffer([]byte(Body))
	req, err := http.NewRequest(Method, Url, body_byte)
	if err != nil {
		return nil
	}

	// Set Requests Headers
	for _, header := range Header {
		grep_key := regexp.MustCompile("(.*): ")
		var header_key = grep_key.FindStringSubmatch(header)[1]
		grep_value := regexp.MustCompile(": (.*)")
		var header_value = grep_value.FindStringSubmatch(header)[1]
		req.Header.Set(header_key, header_value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	// Get Response Body for string
	body_bytes, err := ioutil.ReadAll(resp.Body)
	var response_body = string(body_bytes)
	// Solve the problem of garbled body codes with unmatched numbers
	if !utf8.Valid(body_bytes) {
		data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(body_bytes)
		response_body = string(data)
	}
	// Get Response title
	grep_title := regexp.MustCompile("<title>(.*)</title>")
	if len(grep_title.FindStringSubmatch(response_body)) != 0 {
		resp_title = grep_title.FindStringSubmatch(response_body)[1]
	} else {
		resp_title = "None"
	}
	// Get Response Header for string
	var res []string
	for name, values := range resp.Header {
		for _, value := range values {
			res = append(res, fmt.Sprintf("%s: %s", name, value))
		}
	}
	var response_header string
	for _, re := range res {
		response_header += re + "\n"
	}
	// get response status code
	var status_code = resp.StatusCode
	response_status_code := strconv.Itoa(status_code)
	RespData := []RespLab{
		{Url, response_body, response_header, response_status_code, resp_title, ""},
	}
	return RespData
}

