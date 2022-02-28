package internal

import (
	"fmt"
	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/pkg/logger"
)

func information() {
	if flag.NetWork == flag.InUrl && flag.InUrl == flag.File {
		logger.Fatal(logger.Red("A target must be specified [ -i -u -f ]"))
	}
	if flag.Mode != "" {
		logger.Info(fmt.Sprintf("Discover only with %s protocol", logger.White(flag.Mode)))
	}
	if flag.Type != "" {
		logger.Info(fmt.Sprintf("Discover only with %s", logger.White(flag.Type)))
	}
	if flag.Timeout != 5 {
		logger.Info(fmt.Sprintf("Set the global timeout to %s", logger.White(flag.Timeout)))
	}
	if flag.Thread != 500 && flag.Thread < 1000 {
		logger.Info(fmt.Sprintf("Customize the number of threads to %s", logger.White(flag.Thread)))
	} else if flag.Thread > 999 {
		logger.Warning(fmt.Sprintf("%s %s %s", logger.Yellow("The thread size exceeds"), logger.White(flag.Thread), logger.Yellow("and may be inaccurate")))
	} else {
		logger.Info("The default number of threads is 500")
	}

	if flag.NoIcmp {
		logger.Warning(logger.Yellow("Not use ICMP/PING to detect surviving hosts"))
	}
	if flag.NoColor {
		logger.Warning(logger.Yellow("Don't show color"))
	}
	if flag.Proxy != "" {
		logger.Info(fmt.Sprintf("Scan with proxy %s", logger.White(flag.Proxy)))
	}
}
