package main

/*
   Asset discovery and identification tools
*/

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var url = flag.String("url", "", "Target url.")

func get_time() string {
	var timeStr string = time.Now().Format("15:04:05")
	return timeStr
}



type ColorType int
const (
	ColorType_Invalid ColorType = iota
	ColorType_Green
	ColorType_LightGreen
	ColorType_Cyan
	ColorType_LightCyan
	ColorType_Red
	ColorType_LightRed
	ColorType_Yellow
	ColorType_Black
	ColorType_DarkGray
	ColorType_LightGray
	ColorType_White
	ColorType_Blue
	ColorType_LightBlue
	ColorType_Purple
	ColorType_LightPurple
	ColorType_Brown
)

var (
	ColorFuncMap = map[ColorType]ColorFunc{
		ColorType_Green:       Green,
		ColorType_LightGreen:  LightGreen,
		ColorType_Cyan:        Cyan,
		ColorType_LightCyan:   LightCyan,
		ColorType_Red:         Red,
		ColorType_LightRed:    LightRed,
		ColorType_Yellow:      Yellow,
		ColorType_Black:       Black,
		ColorType_DarkGray:    DarkGray,
		ColorType_LightGray:   LightGray,
		ColorType_White:       White,
		ColorType_Blue:        Blue,
		ColorType_LightBlue:   LightBlue,
		ColorType_Purple:      Purple,
		ColorType_LightPurple: LightPurple,
		ColorType_Brown:       Brown,
	}
	colorFns = []ColorFunc{Green, LightGreen, Cyan, LightCyan, Red, LightRed, Yellow, Black, DarkGray, LightGray, White, Blue, LightBlue, Purple, LightPurple, Brown}
)

//所有的颜色函数
type ColorFunc func(string, ...interface{}) string

//根据颜色获取颜色函数
func GetColorFunc(t ColorType) (ColorFunc, bool) {
	if fn, ok := ColorFuncMap[t]; ok {
		return fn, ok
	}
	return nil, false
}

//绿色字体，modifier里，第一个控制闪烁，第二个控制下划线
func Green(str string, modifier ...interface{}) string {
	return CliColorRender(str, 32, 0, modifier...)
}

//淡绿
func LightGreen(str string, modifier ...interface{}) string {
	return CliColorRender(str, 32, 1, modifier...)
}

//青色/蓝绿色
func Cyan(str string, modifier ...interface{}) string {
	return CliColorRender(str, 36, 0, modifier...)
}

//淡青色
func LightCyan(str string, modifier ...interface{}) string {
	return CliColorRender(str, 36, 1, modifier...)
}

//红字体
func Red(str string, modifier ...interface{}) string {
	return CliColorRender(str, 31, 0, modifier...)
}

//淡红色
func LightRed(str string, modifier ...interface{}) string {
	return CliColorRender(str, 31, 1, modifier...)
}

//黄色字体
func Yellow(str string, modifier ...interface{}) string {
	return CliColorRender(str, 33, 0, modifier...)
}

//黑色
func Black(str string, modifier ...interface{}) string {
	return CliColorRender(str, 30, 0, modifier...)
}

//深灰色
func DarkGray(str string, modifier ...interface{}) string {
	return CliColorRender(str, 30, 1, modifier...)
}

//浅灰色
func LightGray(str string, modifier ...interface{}) string {
	return CliColorRender(str, 37, 0, modifier...)
}

//白色
func White(str string, modifier ...interface{}) string {
	return CliColorRender(str, 37, 1, modifier...)
}

//蓝色
func Blue(str string, modifier ...interface{}) string {
	return CliColorRender(str, 34, 0, modifier...)
}

//淡蓝
func LightBlue(str string, modifier ...interface{}) string {
	return CliColorRender(str, 34, 1, modifier...)
}

//紫色
func Purple(str string, modifier ...interface{}) string {
	return CliColorRender(str, 35, 0, modifier...)
}

//淡紫色
func LightPurple(str string, modifier ...interface{}) string {
	return CliColorRender(str, 35, 1, modifier...)
}

//棕色
func Brown(str string, modifier ...interface{}) string {
	return CliColorRender(str, 33, 0, modifier...)
}

func CliColorRender(str string, color int, weight int, extraArgs ...interface{}) string {
	//闪烁效果
	var isBlink int64 = 0
	if len(extraArgs) > 0 {
		isBlink = reflect.ValueOf(extraArgs[0]).Int()
	}
	//下划线效果
	var isUnderLine int64 = 0
	if len(extraArgs) > 1 {
		isUnderLine = reflect.ValueOf(extraArgs[1]).Int()
	}
	var mo []string
	if isBlink > 0 {
		mo = append(mo, "05")
	}
	if isUnderLine > 0 {
		mo = append(mo, "04")
	}
	if weight > 0 {
		mo = append(mo, fmt.Sprintf("%d", weight))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	buf := bytes.Buffer{}
	buf.WriteString("\033[")
	buf.WriteString(strings.Join(mo, ";"))
	buf.WriteString(";")
	buf.WriteString(fmt.Sprintf("%d", color))
	buf.WriteString("m")
	buf.WriteString(str)
	buf.WriteString("\033[0m")
	//fmt.Sprintf("\033[%s;%dm"+str+"\033[0m", strings.Join(mo, ";"), color)
	return buf.String()
}



func Requests(url string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[!] Error: Could not send the HTTP request.\n")
		os.Exit(1)
	}
	// set requests header
	req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	req.Header.Set("Cookie", "rememberMe=int")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	// get response body for string
	body_bytes, err := ioutil.ReadAll(resp.Body)
	var response_body = string(body_bytes)

	// get response header for string
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
	identify(url, response_body, response_header, response_status_code)
}

func HeaderToArray(header http.Header) (res []string) {
	fmt.Println(header)

	for name, values := range header {
		for _, value := range values {
			res = append(res, fmt.Sprintf("%s: %s", name, value))
			// fmt.Println(res)
		}
	}
	return
}
/*
Rule:
	Name: name
	Type: header, body, url
	Mode: and, or
	Rule
		InBody: str
		InHeader: str
		InUrl: str
*/

type InStr struct {
	InBody string
	InHeader string
	InUrl string
}

type RuleLab struct {
	Name string
	Type string
	Mode string
	Rule InStr
}

var RuleData = []RuleLab{
	{"Apache Shiro", "header", "", InStr{"", "rememberMe=", ""}},
	{"Apache Struts2", "body|url", "or", InStr{"<a href=(.*)\\.action(.*)</a>", "", ".action"}},
}


func identify(url, body string, header string, code string) {
	grep_title := regexp.MustCompile("<title>(.*)</title>")
	title := grep_title.FindStringSubmatch(body)[1]
	var identify_data []string
	for _, rule := range RuleData {
		if rule.Type == "header" {
			if check_heder(url, header, rule.Rule.InHeader, rule.Name, title, code) == true {
				identify_data = append(identify_data, "["+Yellow(rule.Name)+"]")
			}
		}
		if rule.Type == "body" {
			if check_body(url, body, rule.Rule.InBody, rule.Name, title, code) == true {
				identify_data = append(identify_data, "["+Yellow(rule.Name)+"]")
			}
		}
	}
	var identify_result string
	for _, result := range identify_data {
		identify_result += result + " "
	}
	var now_time = Cyan(get_time())
	var succes = LightGreen("+")
	var failed = Red("-")
	// var url = Green(url, 0,1)

	if len(identify_result) != 0 {
		var code = Purple(code)
		var title = Blue(title)
		fmt.Printf("[%s] [%s] [%s] %s%s [%s]\n", now_time, succes, code, identify_result, url, title)
	} else {
		var code = Purple(code)
		var title = Blue(title)
		fmt.Printf("[%s] [%s] [%s] %s [%s]\n", now_time, failed, code, url, title)
	}
	/*
		fmt.Println(Green("字体：Green"))
		fmt.Println(LightGreen("字体：LightGreen"))
		fmt.Println(Cyan("字体：Cyan"))
		fmt.Println(LightCyan("字体：LightCyan"))
		fmt.Println(Red("字体：Red"))
		fmt.Println(LightRed("字体：LightRed"))
		fmt.Println(Yellow("字体：Yellow"))
		fmt.Println(Black("字体：Black"))
		fmt.Println(DarkGray("字体：DarkGray"))
		fmt.Println(LightGray("字体：LightGray"))
		fmt.Println(White("字体：White"))
		fmt.Println(Blue("字体：Blue"))
		fmt.Println(LightBlue("字体：LightBlue"))
		fmt.Println(Purple("字体：Purple"))
		fmt.Println(LightPurple("字体：LightPurple"))
		fmt.Println(Brown("字体：Brown"))
	*/
}

func check_heder(url, response_header string, rule_header string, name string, title string, code string) bool {
	if strings.Index(response_header, rule_header) != -1 {
		return true
	} else {
		return false
	}
}

func check_body(url, response_body string, rule_body string, name string, title string, code string) bool {
	grep := regexp.MustCompile(rule_body)
	if len(grep.FindStringSubmatch(response_body)) != 0 {
		return true
	}

	fmt.Printf("what bug")
	if strings.Index(response_body, rule_body) != -1 {
		return true
	} else {
		return false
	}
}


func main()  {
	flag.Parse()
	fmt.Printf("start..............\n")
	Requests(*url)
}

