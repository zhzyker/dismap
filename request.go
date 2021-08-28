package dismap

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

var defaultHeader = map[string]string{
	"Accept-Language": "zh,zh-TW;q=0.9,en-US;q=0.8,en;q=0.7,zh-CN;q=0.6",
	"User-agent":      "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36",
	"Cookie":          "rememberMe=int",
}

type Sample struct {
	Url        string
	StatusCode int
	Title      string
	Body       string
	Header     string
	Headers    map[string]string
	Cookie     string
	Cookies    map[string]string
	Server     string
	FaviconMd5 string
}

func RequestSample(url string, timeout time.Duration) (*Sample, error) {
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range defaultHeader {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	transformFlag := false
	contentType := resp.Header.Get("Content-Type")
	var reader io.Reader = resp.Body

	if ct := strings.ToUpper(contentType); strings.Contains(ct, "GB2312") || strings.Contains(ct, "GBK") {
		transformFlag = true
		reader = transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	} else if strings.Contains(ct, "BIG5") {
		transformFlag = true
		reader = transform.NewReader(resp.Body, traditionalchinese.Big5.NewDecoder())
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	body := string(b)
	if !transformFlag {
		if !strings.Contains(strings.ToLower(contentType), "utf-8") {
			if lowerBody := strings.ToLower(body); strings.Contains(lowerBody, "gb2312\"") || strings.Contains(lowerBody, "gbk\"") {
				if nBody, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), body); err == nil {
					body = nBody
				}
			}
		}
	}
	// 处理header
	headers := make(map[string]string)
	for key := range resp.Header {
		v := resp.Header.Get(key)
		headers[key] = v
	}
	// 解析cookie
	cookies := make(map[string]string)
	cookie := make([]string, 0)
	for _, v := range resp.Cookies() {
		cookie = append(cookie, v.String())
		cookies[v.Name] = v.Value
	}
	return &Sample{
		Url:        url,
		StatusCode: resp.StatusCode,
		Header:     getHeaderString(resp),
		Headers:    headers,
		Title:      extraTitle(body),
		Body:       body,
		Cookie:     strings.Join(cookie, "\n"),
		Cookies:    cookies,
		Server:     resp.Header.Get("Server"),
		FaviconMd5: faviconMd5(url, client),
	}, nil
}

func faviconMd5(URL string, client *http.Client) string {
	URL, err := joinPath(URL, "/favicon.ico")
	if err != nil {
		return ""
	}
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return ""
	}
	for key, value := range defaultHeader {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	if !(resp.StatusCode == 200) {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum(body))
}

func getHeaderString(resp *http.Response) string {
	r1 := new(http.Response)
	*r1 = *resp
	str, _ := httputil.DumpResponse(r1, false)
	return string(str)
}

var (
	titleTagBegin = `<title>`
	titleTagEnd   = `</title>`
)

func extraTitle(body string) string {
	low := strings.ToLower(body)
	// Some non-utf8 encoding like GB2312 or Big5 may cause ToLower function change the
	// content length, then panic
	if len(low) != len(body) {
		return ""
	}
	begin := strings.Index(low, titleTagBegin)
	if begin < 0 {
		return ""
	}
	begin += len(titleTagBegin)
	end := strings.Index(low, titleTagEnd)
	if end < 0 {
		return ""
	}
	if begin >= end {
		return ""
	}
	return body[begin:end]
}

func joinPath(URL, p string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, p)
	return u.String(), nil
}
