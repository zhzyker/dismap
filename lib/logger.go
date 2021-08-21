package lib

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func get_time() string {
	var timeStr string = time.Now().Format("15:04:05")
	return timeStr
}

func get_timestamp() string {
	var timestampint = time.Now().UnixNano() / 1e6
	timestampstr := strconv.FormatInt(timestampint,10)
	return timestampstr
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
		ColorType_White:       White,
		ColorType_Blue:        Blue,
		ColorType_LightBlue:   LightBlue,
		ColorType_Purple:      Purple,
		ColorType_LightPurple: LightPurple,
		ColorType_Brown:       Brown,
	}
	colorFns = []ColorFunc{Green, LightGreen, Cyan, LightCyan, Red, LightRed, Yellow, White, Blue, LightBlue, Purple, LightPurple, Brown}
)

type ColorFunc func(string, ...interface{}) string

func GetColorFunc(t ColorType) (ColorFunc, bool) {
	if fn, ok := ColorFuncMap[t]; ok {
		return fn, ok
	}
	return nil, false
}

func Green(str string, modifier ...interface{}) string {
	return CliColorRender(str, 32, 0, modifier...)
}

func LightGreen(str string, modifier ...interface{}) string {
	return CliColorRender(str, 32, 1, modifier...)
}

func Cyan(str string, modifier ...interface{}) string {
	return CliColorRender(str, 36, 0, modifier...)
}

func LightCyan(str string, modifier ...interface{}) string {
	return CliColorRender(str, 36, 1, modifier...)
}

func Red(str string, modifier ...interface{}) string {
	return CliColorRender(str, 31, 0, modifier...)
}

func LightRed(str string, modifier ...interface{}) string {
	return CliColorRender(str, 31, 1, modifier...)
}

func Yellow(str string, modifier ...interface{}) string {
	return CliColorRender(str, 33, 1, modifier...)
}

func White(str string, modifier ...interface{}) string {
	return CliColorRender(str, 37, 1, modifier...)
}

func Blue(str string, modifier ...interface{}) string {
	return CliColorRender(str, 34, 0, modifier...)
}

func LightBlue(str string, modifier ...interface{}) string {
	return CliColorRender(str, 34, 1, modifier...)
}

func Purple(str string, modifier ...interface{}) string {
	return CliColorRender(str, 35, 0, modifier...)
}

func LightPurple(str string, modifier ...interface{}) string {
	return CliColorRender(str, 35, 1, modifier...)
}

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


func system() string {
	return runtime.GOOS
}
var sysarch = system()

func logger(level int, log string, detail string) {
	// log level for 0 1 2 3
	// default 0
	if log == "info" {
		if sysarch == "windows" {
			fmt.Printf("[%s] [%s] %s\n", get_time(), "INFO", detail)
		} else {
			fmt.Printf("[%s] [%s] %s\n", Cyan(get_time()), LightGreen("INFO"), detail)
		}
	}
	if log == "warning" {
		if sysarch == "windows" {
			fmt.Printf("[%s] [%s] %s\n", get_time(), "WARNING", detail)
		} else {
			fmt.Printf("[%s] [%s] %s\n", Cyan(get_time()), Yellow("WARNING"), detail)
		}
	}
	if log == "error" {
		if sysarch == "windows" {
			fmt.Printf("[%s] [%s] %s\n", get_time(), "ERROR", detail)
		} else {
			fmt.Printf("[%s] [%s] %s\n", Cyan(get_time()), LightRed("ERROR"), detail)
		}
	}
	if log == "succes" {
		if sysarch == "windows" {
			fmt.Printf("[%s] [%s] %s\n", get_time(), "+", detail)
		} else {
			fmt.Printf("[%s] [%s] %s\n", Cyan(get_time()), LightGreen("+"), detail)
		}
	}
	if log == "failed" {
		if sysarch == "windows" {
			fmt.Printf("[%s] [%s] %s\n", get_time(), "-", detail)
		} else {
			fmt.Printf("[%s] [%s] %s\n", Cyan(get_time()), LightRed("-"), detail)
		}
	}
}
