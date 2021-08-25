package lib

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/zhzyker/dismap/pkg/logger"
)

func ManageFlag() {
	banner()
	flag.Parse()
	Ports := ParsePort(Port)
	//runtime.GOMAXPROCS(4)
	wg := sync.WaitGroup{}
	lock := &sync.Mutex{}
	// output files
	_, err := os.Stat(OutPut)
	if err != nil {
		var dismap_header string
		dismap_header =
			"######          dismap 0.1 output file          ######\r\n" +
				"###### asset discovery and identification tools ######\r\n" +
				"######   by:https://github.com/zhzyker/dismap   ######\r\n"
		f, _ := os.Create(OutPut)
		defer f.Close()
		_, err := f.WriteString(dismap_header)
		if err != nil {
			panic(err)
		}
	}
	fl, _ := os.OpenFile(OutPut, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if NetWork != "" {
		// Start detecting surviving hosts
		logger.Info("Start to detect host from " + NetWork)
		var SurviveHosts []string
		IntAllHost := 0
		IntSurHost := 0
		IntSyncHost := 0
		IntSyncUrl := 0
		hosts, _ := ParseNetHosts(NetWork)
		var ActualHosts []string
		if NoIcmp == false {
			for _, host := range hosts {
				wg.Add(1)
				IntAllHost++
				IntSyncHost++
				go func(host string) {
					if PingHost(host, TimeOut) == true {
						IntSurHost++
						logger.Info("PING found alive host " + host)
						lock.Lock()
						SurviveHosts = append(SurviveHosts, host)
						lock.Unlock()
					}
					wg.Done()
				}(host)
				if IntSyncHost >= Thread {
					IntSyncHost = 0
					wg.Wait()
				}
			}
			wg.Wait()
			ActualHosts = SurviveHosts
			logger.Info(
				logger.LightGreen("There are total of ") +
					logger.White(strconv.Itoa(IntAllHost)) +
					logger.LightGreen(" hosts, and ") +
					logger.White(strconv.Itoa(IntSurHost)) +
					logger.LightGreen(" are surviving"))
			if IntSurHost <= 5 {
				logger.Warn(logger.Yellow("Too few surviving hosts"))
			}
		} else {
			ActualHosts = hosts
			logger.Warn(logger.Yellow("Not use ICMP/PING to detect surviving hosts"))
		}
		logger.Info("Start to identify the targets")
		IntAllUrl := 0
		IntIdeUrl := 0
		for _, host := range ActualHosts {
			for _, port := range Ports {
				wg.Add(1)
				IntSyncUrl++
				url := ParseUrl(host, strconv.Itoa(port))
				go func(url string) {
					var res_type string
					var res_code string
					var res_result string
					var res_result_nc string
					var res_url string
					var res_title string
					for _, results := range Identify(url, TimeOut) {
						res_type = results.Type
						res_code = results.RespCode
						res_result = results.Result
						res_result_nc = results.ResultNc
						res_url = results.Url
						res_title = results.Title
					}
					lock.Lock()
					if len(res_result) != 0 {
						IntIdeUrl++
						IntAllUrl++
						logger.Success("[" + logger.Purple(res_code) + "] " + res_result + res_url + " [" + logger.Blue(res_title) + "]")
						//output(OutPut, lock, "[+] ["+res_code+"] "+ res_result_nc + "{ " + res_url + " } ["+res_title+"]\n")
						content := "[+] [" + res_code + "] " + res_result_nc + "{ " + res_url + " } [" + res_title + "]"
						var text = []byte(content + "\n")
						_, err = fl.Write(text)
						//fmt.Printf("[%s] [%s] [%s] %s%s [%s]\n", now_time, succes, RespCode, identify_result, url, title)
					} else if res_code != "" {
						IntAllUrl++
						logger.Failed("[" + logger.Purple(res_code) + "] " + res_url + " [" + logger.Blue(res_title) + "]")
						//output(OutPut, lock, "[-] ["+res_code+"] " + "{ " + res_url + " } ["+res_title+"]\n")
						content := "[-] [" + res_code + "] " + "{ " + res_url + " } [" + res_title + "]"
						var text = []byte(content + "\n")
						_, err = fl.Write(text)
					}
					lock.Unlock()

					if 1 == 2 { // ahhhhhhhhhhhhhh
						fmt.Println(res_type)
					}
					wg.Done()
				}(url)
				if IntSyncUrl >= Thread {
					IntSyncUrl = 0
					wg.Wait()
				}
			}
		}
		wg.Wait()
		logger.Info(logger.LightGreen("A total of ") +
			logger.White(strconv.Itoa(IntAllUrl)) +
			logger.LightGreen(" urls, the rule base hits ") +
			logger.White(strconv.Itoa(IntIdeUrl)) +
			logger.LightGreen(" urls"))

	} else if Url != "" || Files == "" {
		var res_type string
		var res_code string
		var res_result string
		var res_result_nc string
		var res_url string
		var res_title string
		for _, results := range Identify(Url, TimeOut) {
			res_type = results.Type
			res_code = results.RespCode
			res_result = results.Result
			res_result_nc = results.ResultNc
			res_url = results.Url
			res_title = results.Title
		}
		//lock.Lock()
		if len(res_result) != 0 {
			logger.Success("[" + logger.Purple(res_code) + "] " + res_result + res_url + " [" + logger.Blue(res_title) + "]")
			//output(OutPut, lock, "[+] ["+res_code+"] "+ res_result_nc + "{ " + res_url + " } ["+res_title+"]\n")
			content := "[+] [" + res_code + "] " + res_result_nc + "{ " + res_url + " } [" + res_title + "]"
			var text = []byte(content + "\n")
			_, err = fl.Write(text)
		} else if res_code != "" {
			logger.Failed("[" + logger.Purple(res_code) + "] " + res_url + " [" + logger.Blue(res_title) + "]")
			//output(OutPut, lock, "[-] ["+res_code+"] " + "{ " + res_url + " } ["+res_title+"]\n")
			content := "[-] [" + res_code + "] " + "{ " + res_url + " } [" + res_title + "]"
			var text = []byte(content + "\n")
			_, err = fl.Write(text)
		}
		//lock.Unlock()
		if 1 == 2 { // ahhhhhhhhhhhhhh
			fmt.Println(res_type)
		}
	} else if Url == "" || Files != "" {
		files, err := os.Open(Files)
		if err != nil {
			logger.Error("There is no " + logger.LightRed(Files) + " file or the directory does not exist")
		}
		buf := bufio.NewReader(files)
		IntSyncUrl := 0
		for {
			line, err := buf.ReadString('\n')
			line = strings.TrimSpace(line)
			if err != nil || err == io.EOF {
				break
			}
			if line == "" {
				continue
			}
			IntAllUrl := 0
			IntIdeUrl := 0
			wg.Add(1)
			IntSyncUrl++
			go func(url string) {
				var res_type string
				var res_code string
				var res_result string
				var res_result_nc string
				var res_url string
				var res_title string
				for _, results := range Identify(line, TimeOut) {
					res_type = results.Type
					res_code = results.RespCode
					res_result = results.Result
					res_result_nc = results.ResultNc
					res_url = results.Url
					res_title = results.Title
				}
				lock.Lock()
				if len(res_result) != 0 {
					IntIdeUrl++
					IntAllUrl++
					logger.Success("[" + logger.Purple(res_code) + "] " + res_result + res_url + " [" + logger.Blue(res_title) + "]")
					//output(OutPut, lock, "[+] ["+res_code+"] "+ res_result_nc + "{ " + res_url + " } ["+res_title+"]\n")
					content := "[+] [" + res_code + "] " + res_result_nc + "{ " + res_url + " } [" + res_title + "]"
					var text = []byte(content + "\n")
					_, err = fl.Write(text)
					//fmt.Printf("[%s] [%s] [%s] %s%s [%s]\n", now_time, succes, RespCode, identify_result, url, title)
				} else if res_code != "" {
					IntAllUrl++
					logger.Failed("[" + logger.Purple(res_code) + "] " + res_url + " [" + logger.Blue(res_title) + "]")
					//output(OutPut, lock, "[-] ["+res_code+"] " + "{ " + res_url + " } ["+res_title+"]\n")
					content := "[-] [" + res_code + "] " + "{ " + res_url + " } [" + res_title + "]"
					var text = []byte(content + "\n")
					_, err = fl.Write(text)
				}
				lock.Unlock()
				if 1 == 2 { // ahhhhhhhhhhhhhh
					fmt.Println(res_type)
				}
				wg.Done()
			}(line)
			if IntSyncUrl >= Thread {
				IntSyncUrl = 0
				wg.Wait()
			}
		}
		wg.Wait()
		files.Close()
	}
	fl.Close()
	logger.Info("The identification results are saved in " + OutPut)
	logger.Info("Identification completed and ended")
}

func PingHost(host string, timeout int) bool {
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
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}
