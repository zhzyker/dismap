package lib

import "fmt"

func banner() {
	banners := "_____________\n" +
		"______  /__(_)_____________ _________ ________\n" +
		"_  __  /__  /__  ___/_  __ `__ \\  __ `/__  __ \\\n" +
		"/ /_/ / _  / _(__  )_  / / / / / /_/ /__  /_/ /\n" +
		"\\__,_/  /_/  /____/ /_/ /_/ /_/\\__,_/ _  .___/\n" +
		"                                        /_/"
	infor := "  dismap version: 0.1 release\n" +
		"  author: zhzyker && Nemophllist\n" +
		"  from: https://github.com/zhzyker/dismap\n"
	if sysarch == "windows" {
		fmt.Println(banners)
		fmt.Println(infor)
	} else {
		fmt.Println(LightPurple(banners))
		fmt.Println(White(infor))
	}

}
