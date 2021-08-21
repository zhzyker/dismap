package lib

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
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
		logger(0,"info", "Start to detect host from " + NetWork)
		var SurviveHosts [] string
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
						logger(0,"info", "PING found alive host " + host)
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
			ActualHosts = hosts
			wg.Wait()
			if sysarch == "windows" {
				logger(0,"info", "There are total of "+strconv.Itoa(IntAllHost)+" hosts, and "+strconv.Itoa(IntSurHost)+" are surviving")
			} else {
				logger(0,"info",
					LightGreen("There are total of ")+
						White(strconv.Itoa(IntAllHost))+
						LightGreen(" hosts, and ")+
						White(strconv.Itoa(IntSurHost))+
						LightGreen(" are surviving"))
			}
			if IntSurHost <= 5 {
				if sysarch == "windows" {
					logger(0,"warning", "Too few surviving hosts")
				} else {
					logger(0,"warning", Yellow("Too few surviving hosts"))
				}

			}
		} else {
			ActualHosts = hosts
			if sysarch == "windows" {
				logger(0,"warning", "Not use ICMP/PING to detect surviving hosts")
			} else {
				logger(0,"warning", Yellow("Not use ICMP/PING to detect surviving hosts"))
			}
		}
		logger(0,"info", "Start to identify the targets")
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
						if sysarch == "windows" {
							logger(0,"succes", "["+res_code+"] "+ res_result + res_url + " ["+res_title+"]")
						} else {
							logger(0,"succes", "["+Purple(res_code)+"] "+ res_result + res_url + " ["+Blue(res_title)+"]")
						}
						//output(OutPut, lock, "[+] ["+res_code+"] "+ res_result_nc + "{ " + res_url + " } ["+res_title+"]\n")
						content := "[+] ["+res_code+"] "+ res_result_nc + "{ " + res_url + " } ["+res_title+"]"
						var text = []byte(content + "\n")
						_, err = fl.Write(text)
						//fmt.Printf("[%s] [%s] [%s] %s%s [%s]\n", now_time, succes, RespCode, identify_result, url, title)
					} else if res_code != "" {
						IntAllUrl++
						if sysarch == "windows" {
							logger(0,"failed", "["+res_code+"] " + res_url + " ["+res_title+"]")
						} else {
							logger(0,"failed", "["+Purple(res_code)+"] "+res_url+" ["+Blue(res_title)+"]")
						}
						//output(OutPut, lock, "[-] ["+res_code+"] " + "{ " + res_url + " } ["+res_title+"]\n")
						content := "[-] ["+res_code+"] " + "{ " + res_url + " } ["+res_title+"]"
						var text = []byte(content + "\n")
						_, err = fl.Write(text)
					}
					lock.Unlock()

					if 1==2 { // ahhhhhhhhhhhhhh
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
		if sysarch == "windows" {
			logger(0,"info", "A total of "+strconv.Itoa(IntAllUrl)+" urls, the rule base hits "+strconv.Itoa(IntIdeUrl)+" urls")
		} else {
			logger(0,"info",
				LightGreen("A total of ")+
					White(strconv.Itoa(IntAllUrl))+
					LightGreen(" urls, the rule base hits ")+
					White(strconv.Itoa(IntIdeUrl))+
					LightGreen(" urls"))
		}

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
			if sysarch == "windows" {
				logger(0,"succes", "["+res_code+"] "+ res_result + res_url + " ["+res_title+"]")
			} else {
				logger(0,"succes", "["+Purple(res_code)+"] "+res_result+res_url+" ["+Blue(res_title)+"]")
			}
			//output(OutPut, lock, "[+] ["+res_code+"] "+ res_result_nc + "{ " + res_url + " } ["+res_title+"]\n")
			content := "[+] ["+res_code+"] "+ res_result_nc + "{ " + res_url + " } ["+res_title+"]"
			var text = []byte(content + "\n")
			_, err = fl.Write(text)
		} else if res_code != "" {
			if sysarch == "windows" {
				logger(0,"failed", "["+res_code+"] " + res_url + " ["+res_title+"]")
			} else {
				logger(0,"failed", "["+Purple(res_code)+"] "+res_url+" ["+Blue(res_title)+"]")
			}
			//output(OutPut, lock, "[-] ["+res_code+"] " + "{ " + res_url + " } ["+res_title+"]\n")
			content := "[-] ["+res_code+"] " + "{ " + res_url + " } ["+res_title+"]"
			var text = []byte(content + "\n")
			_, err = fl.Write(text)
		}
		//lock.Unlock()
		if 1==2 { // ahhhhhhhhhhhhhh
			fmt.Println(res_type)
		}
	} else if Url == "" || Files != "" {
		files, err := os.Open(Files)
		if err != nil {
			if sysarch == "windows" {
				logger(0,"error", "There is no " + Files + " file or the directory does not exist")
			} else {
				logger(0,"error", "There is no " + LightRed(Files) + " file or the directory does not exist")
			}

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
					if sysarch == "windows" {
						logger(0,"succes", "["+res_code+"] "+ res_result + res_url + " ["+res_title+"]")
					} else {
						logger(0,"succes", "["+Purple(res_code)+"] "+ res_result + res_url + " ["+Blue(res_title)+"]")
					}
					//output(OutPut, lock, "[+] ["+res_code+"] "+ res_result_nc + "{ " + res_url + " } ["+res_title+"]\n")
					content := "[+] ["+res_code+"] "+ res_result_nc + "{ " + res_url + " } ["+res_title+"]"
					var text = []byte(content + "\n")
					_, err = fl.Write(text)
					//fmt.Printf("[%s] [%s] [%s] %s%s [%s]\n", now_time, succes, RespCode, identify_result, url, title)
				} else if res_code != "" {
					IntAllUrl++
					if sysarch == "windows" {
						logger(0,"failed", "["+res_code+"] " + res_url + " ["+res_title+"]")
					} else {
						logger(0,"failed", "["+Purple(res_code)+"] "+res_url+" ["+Blue(res_title)+"]")
					}
					//output(OutPut, lock, "[-] ["+res_code+"] " + "{ " + res_url + " } ["+res_title+"]\n")
					content := "[-] ["+res_code+"] " + "{ " + res_url + " } ["+res_title+"]"
					var text = []byte(content + "\n")
					_, err = fl.Write(text)
				}
				lock.Unlock()
				if 1==2 { // ahhhhhhhhhhhhhh
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
	logger(0,"info", "The identification results are saved in " + OutPut)
	logger(0,"info", "Identification completed and ended")
}

func PingHost(host string, timeout int) bool {
	var to = strconv.Itoa(timeout)
	var cmd *exec.Cmd
	if sysarch == "windows" {
		cmd = exec.Command("ping", host, "-n", "1", "-w", to)
	} else if sysarch == "linux" {
		cmd = exec.Command("ping", host, "-c", "1", "-w", to, "-W", to)
	} else if sysarch == "darwin" {
		cmd = exec.Command("ping", host, "-c", "1", "-W", to)
	}
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}
