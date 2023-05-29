package logger

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/zhzyker/dismap/internal/flag"

	"github.com/gookit/color"
)

var (
	Red         = color.Red.Render
	Cyan        = color.Cyan.Render
	Yellow      = color.Yellow.Render
	White       = color.White.Render
	Blue        = color.Blue.Render
	Purple      = color.Style{color.Magenta, color.OpBold}.Render
	LightRed    = color.Style{color.Red, color.OpBold}.Render
	LightGreen  = color.Style{color.Green, color.OpBold}.Render
	LightWhite  = color.Style{color.White, color.OpBold}.Render
	LightCyan   = color.Style{color.Cyan, color.OpBold}.Render
	LightYellow = color.Style{color.Yellow, color.OpBold}.Render
	//LightBlue  = color.Style{color.Blue, color.OpBold}.Render
)

var (
	defaultLevel = LevelWarning
)

func SetLevel(l Level) {
	defaultLevel = l
}

func log(l Level, detail string) {
	switch flag.Level {
	case 0:
		SetLevel(0)
	case 1:
		SetLevel(1)
	case 2:
		SetLevel(2)
	case 3:
		SetLevel(3)
	case 4:
		SetLevel(4)
	case 5:
		SetLevel(5)
	}

	if l > defaultLevel {
		return
	}
	if flag.NoColor {
		fmt.Println(Clean(detail))
		return
	} else {
		fmt.Println(detail)
	}
	if l == LevelFatal {
		os.Exit(0)
	}
}

func Fatal(detail string) {
	log(LevelFatal, fmt.Sprintf("[%s] [%s] %s", Cyan(GetTime()), LightRed("FATAL"), detail))
}

func Error(detail string) {
	log(LevelError, fmt.Sprintf("[%s] [%s] %s", Cyan(GetTime()), LightRed("ERROR"), detail))
}

func Info(detail string) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(GetTime()), LightGreen("INFO"), detail))
}

func Warning(detail string) {
	log(LevelWarning, fmt.Sprintf("[%s] [%s] %s", Cyan(GetTime()), LightYellow("WARNING"), detail))
}

func Debug(detail string) {
	log(LevelDebug, fmt.Sprintf("[%s] [%s] %s", Cyan(GetTime()), LightWhite("DEBUG"), detail))
}

func Verbose(detail string) {
	log(LevelVerbose, fmt.Sprintf("[%s] [%s] %s", Cyan(GetTime()), LightCyan("VERBOSE"), detail))
}

func Success(detail string) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(GetTime()), LightGreen("+"), detail))
}

func Failed(detail string) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(GetTime()), LightRed("-"), detail))
}

func GetTime() string {
	return time.Now().Format("15:04:05")
}

func DebugError(err error) bool {
	/* Processing error display */
	if err != nil {
		pc, _, line, _ := runtime.Caller(1)
		Debug(fmt.Sprintf("%s%s%s",
			White(runtime.FuncForPC(pc).Name()),
			LightWhite(fmt.Sprintf(" line:%d ", line)),
			White(err)))
		return true
	}
	return false
}

// Clean by https://github.com/acarl005/stripansi/blob/master/stripansi.go
func Clean(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}
