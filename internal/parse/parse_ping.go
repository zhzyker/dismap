package parse

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"os/exec"
	"runtime"
	"strconv"
)

func Ping(host string, timeout int) bool {
	var to = strconv.Itoa(timeout)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", host, "-n", "1", "-w", to)
	case "linux":
		cmd = exec.Command("ping", host, "-c", "1", "-w", to, "-W", to)
	case "darwin":
		cmd = exec.Command("ping", host, "-c", "1", "-W", to)
	}
	if cmd == nil {
		return false
	}
	err := cmd.Run()
	if logger.DebugError(err) {
		return false
	}
	return true
}
