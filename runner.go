package dismap

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zhzyker/dismap/pkg/limiter"
	"github.com/zhzyker/dismap/pkg/logger"
)

var (
	outputHeader = `######          dismap 0.1 output file          ######
###### asset discovery and identification tools ######
######   by:https://github.com/zhzyker/dismap   ######
`
)

type ScanResult struct {
	Sample       *Sample
	Fingerprints []string
}

type Runner struct {
	options  *Options
	limiter  *limiter.Limiter
	targetCh chan string
	outputCh chan ScanResult
}

func NewRunner(p *Options) *Runner {
	return &Runner{
		options: p,
		limiter: limiter.New(int64(p.Threads)),
	}
}

func (r *Runner) Scan() error {
	wg := sync.WaitGroup{}
	r.targetCh = make(chan string, 50)
	r.outputCh = make(chan ScanResult, 50)

	if err := r.mergeTargets(); err != nil {
		return err
	}

	file, err := r.parseOutputFile()
	if err != nil {
		return err
	}
	defer file.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		r.output(file)
	}()

	logger.Infoln("Start to identify the targets")
	timeout := time.Duration(r.options.TimeOut) * time.Second
	for {
		select {
		case url, ok := <-r.targetCh:
			if url == "" {
				if ok {
					continue
				}
				r.targetCh = nil
				goto LOOP_BREAK
			}
			r.limiter.Allow()
			go func(url string) {
				defer r.limiter.Done()
				req, err := MakeDefaultRequest(url)
				if err != nil {
					return
				}
				sample, err := RequestSample(req, timeout)
				if err != nil {
					return
				}
				fingerprints := IdentifyRules(sample, timeout)
				r.outputCh <- ScanResult{
					Sample:       sample,
					Fingerprints: fingerprints,
				}
			}(url)
		}
	LOOP_BREAK:
		if r.targetCh == nil {
			break
		}
	}

	r.limiter.Wait()
	close(r.outputCh)
	wg.Wait()

	logger.Infoln("The identification results are saved in", r.options.OutPut)
	logger.Infoln("Identification completed and ended")
	return nil
}

func (r *Runner) output(file *os.File) {
	for {
		select {
		case res, ok := <-r.outputCh:
			if res.Sample == nil {
				if ok {
					continue
				}
				r.outputCh = nil
				goto LOOP_BREAK
			}
			var content string
			if len(res.Fingerprints) > 0 {
				fmtfps := strings.Join(res.Fingerprints, "] [")
				fmtfps = "[" + fmtfps + "]"
				content = fmt.Sprintf("[+] [%d] %v %s [%s]\n", res.Sample.StatusCode, fmtfps, res.Sample.Url, res.Sample.Title)
				logger.Successf("[%s] %v %s [%s]\n", logger.Magenta(res.Sample.StatusCode), fmtfps, res.Sample.Url, logger.Blue(res.Sample.Title))
			} else {
				content = fmt.Sprintf("[-] [%d] %s [%s]\n", res.Sample.StatusCode, res.Sample.Url, res.Sample.Title)
				logger.Failedf("[%s] %s [%s]\n", logger.Magenta(res.Sample.StatusCode), res.Sample.Url, logger.Blue(res.Sample.Title))
			}
			file.WriteString(content)
		}
	LOOP_BREAK:
		if r.outputCh == nil {
			break
		}
	}
}

func (r *Runner) parseOutputFile() (*os.File, error) {
	var (
		file *os.File
		err  error
		init bool
	)
	_, err = os.Stat(r.options.OutPut)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		init = true
	}

	file, err = os.OpenFile(r.options.OutPut, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	if init {
		file.WriteString(outputHeader)
	}
	return file, nil
}

func (r *Runner) mergeTargets() error {
	lock := sync.Mutex{}

	var (
		ips []string
		err error
	)

	if r.options.Ips != "" {
		logger.Infoln("Start to detect host from", r.options.Ips)
		ips, err = ParseIPRange(r.options.Ips)
		if err != nil {
			return err
		}
	}

	domains := make([]string, 0, 1)
	urls := make([]string, 0, 1)
	if r.options.Url != "" && isURL(r.options.Url) {
		urls = append(urls, r.options.Url)
	}

	if r.options.File != "" {
		f, err := os.Open(r.options.File)
		if err != nil {
			logger.Errorln("There is no " + logger.LightRed(r.options.File) + " file or the directory does not exist")
		} else {
			defer f.Close()
			s := bufio.NewScanner(f)
			for s.Scan() {
				t := strings.TrimSpace(s.Text())
				if strings.HasPrefix(t, "http") && isURL(t) {
					urls = append(urls, t)
					continue
				}
				if isDomainName(t) {
					domains = append(domains, t)
					continue
				}
				if strings.Contains(t, ":") {
					parts := strings.Split(t, ":")
					if len(parts) != 2 {
						continue
					}
					port, err := strconv.Atoi(parts[0])
					if err != nil {
						continue
					}
					h := strings.TrimSpace(parts[1])
					if isIP(h) && (port > 0 && port < 65535) {
						urls = append(urls, ParseUrl(h, port))
					}
				} else if isIP(t) {
					ips = append(ips, t)
				}
			}
		}
	}

	// 探测存活IP
	hosts := make([]string, 0, len(ips))
	if !r.options.NoIcmp && len(ips) > 0 {
		sum := 0
		survives := make([]string, 0, len(ips))

		for i := range ips {
			r.limiter.Allow()
			go func(h string) {
				defer r.limiter.Done()
				if pingHost(h, r.options.TimeOut) {
					lock.Lock()
					sum++
					logger.Infoln("PING found alive host", h)
					survives = append(survives, h)
					lock.Unlock()
				}
			}(ips[i])
		}

		r.limiter.Wait()
		hosts = survives
		logger.Infof("%s %s %s %s %s\n", logger.LightGreen("There are total of"), logger.White(strconv.Itoa(len(ips))), logger.LightGreen("hosts, and"), logger.White(strconv.Itoa(sum)), logger.LightGreen("are surviving"))
		if sum <= 5 {
			logger.Warnln(logger.Yellow("Too few surviving hosts"))
		}
	} else {
		hosts = ips
		logger.Warnln(logger.Yellow("Not use ICMP/PING to detect surviving hosts"))
	}

	go func() {
		defer close(r.targetCh)
		if len(hosts) > 0 {
			if ports, err := ParsePorts(r.options.Ports); err == nil {
				for i := range hosts {
					for j := range ports {
						r.targetCh <- ParseUrl(hosts[i], ports[j])
					}
				}
			}
		}
		for i := range urls {
			r.targetCh <- urls[i]
		}
		for i := range domains {
			r.targetCh <- "http://" + domains[i]
		}
	}()
	return nil
}

func pingHost(host string, timeout int) bool {
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
	return cmd.Run() == nil
}

func isURL(r string) bool {
	_, err := url.Parse(r)
	return err == nil
}

func isIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

var (
	domainRex = regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`)
)

func isDomainName(d string) bool {
	return domainRex.MatchString(d)
}
