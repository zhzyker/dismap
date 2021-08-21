package lib

import (
	"os"
	"sync"
)


// Multithreading bug, opening files too many times
func output(filename string, lock *sync.Mutex, content string) {
	_, err := os.Stat(filename)
	if err == nil {
		var text = []byte(content + "\n")
		lock.Lock()
		fl, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			//fmt.Printf("Open %s error, %v\n", filename, err)
		}
		_, err = fl.Write(text)
		fl.Close()
		if err != nil {
			//fmt.Printf("Write %s error, %v\n", filename, err)
		}
		defer lock.Unlock()
	} else {
		var dismap_header string
		dismap_header =
				"######          dismap 0.1 output file          ######\r\n" +
				"###### asset discovery and identification tools ######\r\n" +
				"######   by:https://github.com/zhzyker/dismap   ######\r\n"
		f, _ := os.Create(filename)
		defer f.Close()
		_, err := f.WriteString(dismap_header + content)
		if err != nil {
			panic(err)
		}
	}
}
