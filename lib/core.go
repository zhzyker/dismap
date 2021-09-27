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
	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}
	// output files
	_, err := os.Stat(OutPut)
	if err != nil {
		var dismap_header string
		dismap_header =
			"######          dismap 0.2 output file          ######\r\n" +
				"###### asset discovery and identification tools ######\r\n" +
				"######   by:https://github.com/zhzyker/dismap   ######\r\n"
		f, _ := os.Create(OutPut)
		defer f.Close()
		_, err := f.WriteString(dismap_header)
		if err != nil {
			panic(err)
		}
	}

	if NetWork != "" {
		TargetNetwork(wg, lock, Ports, NetWork)

	} else if InUrl != "" || Files == "" {
		if_url, err := JudgeUrl(InUrl)
		if err == nil {
			TargetUrl(wg, lock, Ports, if_url)
		}
	} else if InUrl == "" || Files != "" {
		files, err := os.Open(Files)
		if err != nil {
			logger.Error("There is no " + logger.LightRed(Files) + " file or the directory does not exist")
		} else {
			if Thread == 508 {
				logger.Info("The default number of threads is 500")
			}
			logger.Info(logger.LightGreen("Batch scan the targets in " + logger.Yellow(Files) + logger.LightGreen(", priority network segment")))
		}
		buf := bufio.NewReader(files)
		var urls []string
		for {
			line, err := buf.ReadString('\n')
			line = strings.TrimSpace(line)
			if err != nil || err == io.EOF {
				break
			}
			if line == "" {
				continue
			}
			_, err = JudgeNet(line)
			if err == nil {
				TargetNetwork(wg, lock, Ports, line)
				continue
			}
			if_url, err := JudgeUrl(line)
			if err == nil {
				urls = append(urls, if_url)
			} else {
				logger.Warning(logger.Yellow(line) + " is not a legal url, please check")
			}
		}
		logger.Info(logger.LightGreen("Start batch identify urls"))
		IntSync := 0
		IntAll := 0
		for _, target := range urls {
			IntSync++
			IntAll++
			wg.Add(1)
			go func(target string) {
				lock.Lock()
				TargetUrl(wg, lock, Ports, target)
				lock.Unlock()
				wg.Done()
			}(target)
			if IntSync >= Thread {
				IntSync = 0
				wg.Wait()
			}
		}
		wg.Wait()
		logger.Info(
			logger.LightGreen("A total of ") +
			logger.LightWhite(strconv.Itoa(IntAll)) +
			logger.LightGreen(" url targets"))
		files.Close()
	}
	logger.Info("The identification results are saved in " + logger.Yellow(OutPut))
	logger.Info("Identification completed and ended")
}

func TargetNetwork(wg *sync.WaitGroup, lock *sync.Mutex, Ports []int, Targets string) {
	fl, _ :=  os.OpenFile(OutPut, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// Start detecting surviving hosts
	logger.Info("Start to detect host from " + Targets)
	var SurviveHosts []string
	IntAllHost := 0
	IntSurHost := 0
	IntSyncHost := 0
	IntSyncUrl := 0
	hosts, _ := ParseNetHosts(Targets)
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
			logger.Warning(logger.Yellow("Too few surviving hosts"))
		}
	} else {
		ActualHosts = hosts
		logger.Warning(logger.Yellow("Not use ICMP/PING to detect surviving hosts"))
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
					content := "[+] [" + res_code + "] " + res_result_nc + "{ " + res_url + " } [" + res_title + "]"
					var text = []byte(content + "\n")
					fl.Write(text)
				} else if res_code != "" {
					IntAllUrl++
					logger.Failed("[" + logger.Purple(res_code) + "] " + res_url + " [" + logger.Blue(res_title) + "]")
					content := "[-] [" + res_code + "] " + "{ " + res_url + " } [" + res_title + "]"
					var text = []byte(content + "\n")
					fl.Write(text)
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
	fl.Close()
}

func TargetUrl(wg *sync.WaitGroup, lock *sync.Mutex, Ports []int, Targets string) {
	fl, _ :=  os.OpenFile(OutPut, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	var res_type string
	var res_code string
	var res_result string
	var res_result_nc string
	var res_url string
	var res_title string
	for _, results := range Identify(Targets, TimeOut) {
		res_type = results.Type
		res_code = results.RespCode
		res_result = results.Result
		res_result_nc = results.ResultNc
		res_url = results.Url
		res_title = results.Title
	}
	if len(res_result) != 0 {
		logger.Success("[" + logger.Purple(res_code) + "] " + res_result + res_url + " [" + logger.Blue(res_title) + "]")
		content := "[+] [" + res_code + "] " + res_result_nc + "{ " + res_url + " } [" + res_title + "]"
		var text = []byte(content + "\n")
		fl.Write(text)
	} else if res_code != "" {
		logger.Failed("[" + logger.Purple(res_code) + "] " + res_url + " [" + logger.Blue(res_title) + "]")
		content := "[-] [" + res_code + "] " + "{ " + res_url + " } [" + res_title + "]"
		var text = []byte(content + "\n")
		fl.Write(text)
	}
	if 1 == 2 { // ahhhhhhhhhhhhhh
		fmt.Println(res_type)
	}
	fl.Close()
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
