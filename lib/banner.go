package lib

import (
	"fmt"

	"github.com/zhzyker/dismap/pkg/logger"
)

func banner() {
	banners := "_____________\n" +
		"______  /__(_)_____________ _________ ________\n" +
		"_  __  /__  /__  ___/_  __ `__ \\  __ `/__  __ \\\n" +
		"/ /_/ / _  / _(__  )_  / / / / / /_/ /__  /_/ /\n" +
		"\\__,_/  /_/  /____/ /_/ /_/ /_/\\__,_/ _  .___/\n" +
		"                                        /_/"
	infor := "  dismap version: 0.2 release\n" +
		"  author: zhzyker && Nemophllist\n" +
		"  from: https://github.com/zhzyker/dismap\n"
	fmt.Println(logger.Purple(banners))
	fmt.Println(logger.LightWhite(infor))
}
