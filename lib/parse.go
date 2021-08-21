package lib

import (
	"github.com/zhzyker/dismap/config"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
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
		if sysarch == "windows" {
			logger(0,"error", cidr+" is not a valid network segment")
			logger(0,"error", "The host format is incorrect, there is no detectable host")
		} else {
			logger(0,"error", LightRed(cidr)+" is not a valid network segment")
			logger(0,"error", LightRed("The hosts format is incorrect, there is no detectable hosts"))
		}
		os.Exit(0)
	}
	if len(regexp.MustCompile("-").FindAllStringIndex(cidr, -1)) == 1 {
		grep := regexp.MustCompile("(.*)-(.*)")
		cidr_ip := grep.FindStringSubmatch(cidr)[1]
		cidr_end := grep.FindStringSubmatch(cidr)[2]
		iperr := net.ParseIP(cidr_ip)
		if iperr == nil {
			if sysarch == "windows" {
				logger(0,"error", cidr+" is not a valid network segment")
				logger(0,"error", "The host format is incorrect, there is no detectable host")
			} else {
				logger(0,"error", LightRed(cidr)+" is not a valid network segment")
				logger(0,"error", LightRed("The hosts format is incorrect, there is no detectable hosts"))
			}
			os.Exit(0)
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
				if sysarch == "windows" {
					logger(0,"error", cidr+" is not a valid network segment")
				} else {
					logger(0,"error", LightRed(cidr)+" is not a valid network segment")
				}
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
				if sysarch == "windows" {
					logger(0,"error", cidr+" is not a valid network segment")
				} else {
					logger(0,"error", LightRed(cidr)+" is not a valid network segment")
				}
				return ips, nil
			}
			for i := int_ip_start; i <= int_ip_end; i++ {
				ip := str_ip_net + strconv.Itoa(i)
				ips = append(ips, ip)
			}
			return ips, nil
		} else {
			if sysarch == "windows" {
				logger(0,"error", cidr+" is not a valid network segment")
				logger(0,"error", "The host format is incorrect, there is no detectable host")
			} else {
				logger(0,"error", LightRed(cidr)+" is not a valid network segment")
				logger(0,"error", LightRed("The hosts format is incorrect, there is no detectable hosts"))
			}
			os.Exit(0)
		}
	} else if len(regexp.MustCompile("/").FindAllStringIndex(cidr, -1)) == 1 {
		ip, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			if sysarch == "windows" {
				logger(0,"error", cidr+" is not a valid network segment")
			} else {
				logger(0,"error", LightRed(cidr)+" is not a valid network segment")
			}
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
				if sysarch == "windows" {
					logger(0,"error", check_host+" is not a valid network segment")
				} else {
					logger(0,"error", LightRed(check_host)+" is not a valid network segment")
				}
			} else {
				i++
			}
		}
		if i == 0 {
			if sysarch == "windows" {
				logger(0,"error", "The host format is incorrect, there is no detectable host")
			} else {
				logger(0,"error", LightRed("The hosts format is incorrect, there is no detectable hosts"))
			}
			os.Exit(0)
		}

		return CustomPorts, nil
	}
	return ips, nil
}

func ParseUrl(host string, port string) string {
	if port == "80" {
		return "http://" + host
	} else if port == "443" {
		return  "https://" + host
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
			p,_ := strconv.Atoi(port)
			ports = append(ports, p)
		}
		return ports
	}
}
