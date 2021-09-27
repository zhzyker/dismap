package lib

import (
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/zhzyker/dismap/config"
	"github.com/zhzyker/dismap/pkg/logger"
)

func incc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		//fmt.Println(ip[j])
		if ip[j] > 0 {
			break
		}
	}
}

func ParseNetHosts(cidr string) ([]string, error) {
	var ips []string
	if len(regexp.MustCompile("0.0.0.0").FindAllStringIndex(cidr, -1)) == 1 {
		logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
		logger.Fatal(logger.LightRed("The hosts format is incorrect, there is no detectable hosts"))
	}
	if len(regexp.MustCompile("-").FindAllStringIndex(cidr, -1)) == 1 {
		grep := regexp.MustCompile("(.*)-(.*)")
		cidr_ip := grep.FindStringSubmatch(cidr)[1]
		cidr_end := grep.FindStringSubmatch(cidr)[2]
		iperr := net.ParseIP(cidr_ip)
		if iperr == nil {
			logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
			logger.Fatal(logger.LightRed("The hosts format is incorrect, there is no detectable hosts"))
		}
		enderr := net.ParseIP(cidr_end)
		int_end, _ := strconv.Atoi(cidr_end)
		// example: 192.168.1.1-1
		if enderr == nil && int_end >= 1 && int_end <= 255 {
			grep := regexp.MustCompile("(.*)\\.(.*)\\.(.*)\\.(.*)")
			str_ip_net := grep.FindStringSubmatch(cidr_ip)[1] + "." + grep.FindStringSubmatch(cidr_ip)[2] + "." + grep.FindStringSubmatch(cidr_ip)[3] + "."
			cidr_start := grep.FindStringSubmatch(cidr_ip)[4]
			int_ip_start, _ := strconv.Atoi(cidr_start)
			int_ip_end, _ := strconv.Atoi(cidr_end)
			if int_ip_start == 0 || int_ip_end >= 255 {
				logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
				return ips, nil
			}
			for i := int_ip_start; i <= int_ip_end; i++ {
				ip := str_ip_net + strconv.Itoa(i)
				ips = append(ips, ip)
			}
			return ips, nil
		} else if enderr != nil && int_end == 0 { // example: 192.168.1.1-192.168.1.10
			grep := regexp.MustCompile("(.*)\\.(.*)\\.(.*)\\.(.*)")
			str_ip_net := grep.FindStringSubmatch(cidr_ip)[1] + "." + grep.FindStringSubmatch(cidr_ip)[2] + "." + grep.FindStringSubmatch(cidr_ip)[3] + "."
			cidr_start := grep.FindStringSubmatch(cidr_ip)[4]
			int_ip_start, _ := strconv.Atoi(cidr_start)
			cidr_end := grep.FindStringSubmatch(cidr_end)[4]
			int_ip_end, _ := strconv.Atoi(cidr_end)
			if int_ip_start == 0 || int_ip_end >= 255 {
				logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
				return ips, nil
			}
			for i := int_ip_start; i <= int_ip_end; i++ {
				ip := str_ip_net + strconv.Itoa(i)
				ips = append(ips, ip)
			}
			return ips, nil
		} else {
			logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
			logger.Fatal(logger.LightRed("The hosts format is incorrect, there is no detectable hosts"))
		}
	} else if len(regexp.MustCompile("/").FindAllStringIndex(cidr, -1)) == 1 {
		ip, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
			return nil, err
		}
		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incc(ip) {
			ips = append(ips, ip.String())
		}
		return ips[1 : len(ips)-1], nil
	} else {
		CustomPorts := strings.Split(cidr, ",")
		i := 0
		for _, check_host := range CustomPorts {
			err := net.ParseIP(check_host)
			if err == nil {
				logger.Error(logger.LightRed(check_host) + " is not a valid network segment")
			} else {
				i++
			}
		}
		if i == 0 {
			logger.Fatal(logger.LightRed("The hosts format is incorrect, there is no detectable hosts"))
		}

		return CustomPorts, nil
	}
	return ips, nil
}

func ParseUrl(host string, port string) string {
	if port == "80" {
		return "http://" + host
	} else if port == "443" {
		return "https://" + host
	} else if len(regexp.MustCompile("443").FindAllStringIndex(port, -1)) == 1 {
		return "https://" + host + ":" + port
	} else {
		return "http://" + host + ":" + port
	}

}

func ParsePort(portstr string) []int {
	var ports []int
	if portstr == "" {
		defport := config.DefaultPorts
		return defport
	} else if len(regexp.MustCompile("-").FindAllStringIndex(portstr, -1)) == 1 {
		grep := regexp.MustCompile("(.*)-(.*)")
		port_start := grep.FindStringSubmatch(portstr)[1]
		port_end := grep.FindStringSubmatch(portstr)[2]
		int_port_start, _ := strconv.Atoi(port_start)
		int_port_end, _ := strconv.Atoi(port_end)
		for p := int_port_start; p <= int_port_end; p++ {
			ports = append(ports, p)

		}
		return ports
	} else {
		CustomPorts := strings.Split(portstr, ",")
		for _, port := range CustomPorts {
			p, _ := strconv.Atoi(port)
			ports = append(ports, p)
		}
		return ports
	}
}

func JudgeUrl(target string) (string, error) {
	re, err := regexp.Compile("(http|https):\\/\\/[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-\\.,@?^=%&:/~\\+#]*[\\w\\-\\@?^=%&/~\\+#])?")
	result := re.FindAllStringSubmatch(target, -1)
	if err != nil || result == nil {
		re := regexp.MustCompile("(.*)\\.(.*)")
		result := re.FindAllStringSubmatch(target, -1)
		if result != nil {
			if len(regexp.MustCompile(":443").FindAllStringIndex(target, -1)) == 1 {
				re := regexp.MustCompile("https:\\/\\/[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-\\.,@?^=%&:/~\\+#]*[\\w\\-\\@?^=%&/~\\+#])?")
				result := re.FindAllStringSubmatch(target, -1)
				if result == nil {
					return "https://" + target, nil
				}
			} else {
				re := regexp.MustCompile("http:\\/\\/[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-\\.,@?^=%&:/~\\+#]*[\\w\\-\\@?^=%&/~\\+#])?")
				result := re.FindAllStringSubmatch(target, -1)
				if result == nil {
					return "http://" + target, nil
				}
			}
		}
		_, err := url.ParseRequestURI(target)
		if err != nil {
			return target, err
		}
	}
	return target, nil
}

func JudgeNet(target string) ([]string, error) {
	var ips []string
	error_ip := "1.1"
	if len(regexp.MustCompile("-").FindAllStringIndex(target, -1)) == 1 {
		grep := regexp.MustCompile("(.*)-(.*)")
		cidr_ip := grep.FindStringSubmatch(target)[1]
		cidr_end := grep.FindStringSubmatch(target)[2]
		iperr := net.ParseIP(cidr_ip)
		if iperr == nil {
			_, _, err := net.ParseCIDR(error_ip)
			if err != nil {
				return ips, err
			}
		}
		enderr := net.ParseIP(cidr_end)
		int_end, _ := strconv.Atoi(cidr_end)
		// example: 192.168.1.1-1
		if enderr == nil && int_end >= 1 && int_end <= 255 {
			grep := regexp.MustCompile("(.*)\\.(.*)\\.(.*)\\.(.*)")
			str_ip_net := grep.FindStringSubmatch(cidr_ip)[1] + "." + grep.FindStringSubmatch(cidr_ip)[2] + "." + grep.FindStringSubmatch(cidr_ip)[3] + "."
			cidr_start := grep.FindStringSubmatch(cidr_ip)[4]
			int_ip_start, _ := strconv.Atoi(cidr_start)
			int_ip_end, _ := strconv.Atoi(cidr_end)
			if int_ip_start == 0 || int_ip_end >= 255 {
				_, _, err := net.ParseCIDR(error_ip)
				if err != nil {
					return ips, err
				}
			}
			for i := int_ip_start; i <= int_ip_end; i++ {
				ip := str_ip_net + strconv.Itoa(i)
				ips = append(ips, ip)
			}
			return ips, nil
		} else if enderr != nil && int_end == 0 { // example: 192.168.1.1-192.168.1.10
			grep := regexp.MustCompile("(.*)\\.(.*)\\.(.*)\\.(.*)")
			str_ip_net := grep.FindStringSubmatch(cidr_ip)[1] + "." + grep.FindStringSubmatch(cidr_ip)[2] + "." + grep.FindStringSubmatch(cidr_ip)[3] + "."
			cidr_start := grep.FindStringSubmatch(cidr_ip)[4]
			int_ip_start, _ := strconv.Atoi(cidr_start)
			cidr_end := grep.FindStringSubmatch(cidr_end)[4]
			int_ip_end, _ := strconv.Atoi(cidr_end)
			if int_ip_start == 0 || int_ip_end >= 255 {
				_, _, err := net.ParseCIDR(error_ip)
				if err != nil {
					return ips, err
				}
			}
			for i := int_ip_start; i <= int_ip_end; i++ {
				ip := str_ip_net + strconv.Itoa(i)
				ips = append(ips, ip)
			}
			return ips, nil
		} else {
			_, _, err := net.ParseCIDR(error_ip)
			if err != nil {
				return ips, err
			}
		}
	} else if len(regexp.MustCompile("/").FindAllStringIndex(target, -1)) == 1 {
		ip, ipnet, err := net.ParseCIDR(target)
		if err != nil {
			return nil, err
		}
		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incc(ip) {
			ips = append(ips, ip.String())
		}
		return ips[1 : len(ips)-1], nil
	} else {
		CustomHosts := strings.Split(target, ",")
		i := 0
		for _, check_host := range CustomHosts {
			_, _, err := net.ParseCIDR(check_host)
			if err != nil {
				return ips, err
			} else {
				i++
			}
		}
		return CustomHosts, nil
	}
	return ips, nil
}
