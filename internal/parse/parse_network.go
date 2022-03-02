package parse

import (
	"github.com/zhzyker/dismap/pkg/logger"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func inCC(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func NetworkParse(cidr string) ([]string, error) {
	var ips []string
	if len(regexp.MustCompile("0.0.0.0").FindAllStringIndex(cidr, -1)) == 1 {
		logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
		logger.Fatal(logger.LightRed("The hosts format is incorrect, there is no detectable hosts"))
	}
	if len(regexp.MustCompile("-").FindAllStringIndex(cidr, -1)) == 1 {
		grep := regexp.MustCompile("(.*)-(.*)")
		cidrIp := grep.FindStringSubmatch(cidr)[1]
		cidrEnd := grep.FindStringSubmatch(cidr)[2]
		ipErr := net.ParseIP(cidrIp)
		if ipErr == nil {
			logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
			logger.Fatal(logger.LightRed("The hosts format is incorrect, there is no detectable hosts"))
		}
		endErr := net.ParseIP(cidrEnd)
		intEnd, err := strconv.Atoi(cidrEnd)
		// example: 192.168.1.1-1
		if endErr == nil && intEnd >= 1 && intEnd <= 255 {
			grep := regexp.MustCompile("(.*)\\.(.*)\\.(.*)\\.(.*)")
			strIpNet := grep.FindStringSubmatch(cidrIp)[1] + "." + grep.FindStringSubmatch(cidrIp)[2] + "." + grep.FindStringSubmatch(cidrIp)[3] + "."
			cidrStart := grep.FindStringSubmatch(cidrIp)[4]
			intIpStart, _ := strconv.Atoi(cidrStart)
			intIpEnd, _ := strconv.Atoi(cidrEnd)
			if intIpStart == 0 || intIpEnd >= 255 {
				logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
				return ips, err
			}
			for i := intIpStart; i <= intIpEnd; i++ {
				ip := strIpNet + strconv.Itoa(i)
				ips = append(ips, ip)
			}
			return ips, err
		} else if endErr != nil && intEnd == 0 { // example: 192.168.1.1-192.168.1.10
			grep := regexp.MustCompile("(.*)\\.(.*)\\.(.*)\\.(.*)")
			strIpNet := grep.FindStringSubmatch(cidrIp)[1] + "." + grep.FindStringSubmatch(cidrIp)[2] + "." + grep.FindStringSubmatch(cidrIp)[3] + "."
			cidrStart := grep.FindStringSubmatch(cidrIp)[4]
			intIpStart, _ := strconv.Atoi(cidrStart)
			cidrEnd := grep.FindStringSubmatch(cidrEnd)[4]
			intIpEnd, err := strconv.Atoi(cidrEnd)
			if intIpStart == 0 || intIpEnd >= 255 {
				logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
				return ips, err
			}
			for i := intIpStart; i <= intIpEnd; i++ {
				ip := strIpNet + strconv.Itoa(i)
				ips = append(ips, ip)
			}
			return ips, err
		} else {
			logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
			logger.Fatal(logger.LightRed("The hosts format is incorrect, there is no detectable hosts"))
		}
	} else if len(regexp.MustCompile("/").FindAllStringIndex(cidr, -1)) == 1 {
		ip, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			logger.Error(logger.LightRed(cidr) + " is not a valid network segment")
			return nil, err
		}
		for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inCC(ip) {
			ips = append(ips, ip.String())
		}
		r := ips[1:len(ips)-1]
		return r, nil
	} else {
		customPorts := strings.Split(cidr, ",")
		i := 0
		for _, checkHost := range customPorts {
			err := net.ParseIP(checkHost)
			if err == nil {
				logger.Error(logger.LightRed(checkHost) + " is not a valid network segment")
			} else {
				i++
			}
		}
		if i == 0 {
			logger.Fatal(logger.LightRed("The hosts format is incorrect, there is no detectable hosts"))
		}
		return customPorts, nil
	}
	return ips, nil
}


func NetJudgeParse(target string) bool {
	address := net.ParseIP(target)
	if address != nil {
		return true
	}

	errorIp := "1.1"
	if len(regexp.MustCompile("-").FindAllStringIndex(target, -1)) == 1 {
		grep := regexp.MustCompile("(.*)-(.*)")
		cidrIp := grep.FindStringSubmatch(target)[1]
		cidrEnd := grep.FindStringSubmatch(target)[2]
		ipErr := net.ParseIP(cidrIp)
		if ipErr == nil {
			_, _, err := net.ParseCIDR(errorIp)
			if err != nil {
				return false
			}
		}
		endErr := net.ParseIP(cidrEnd)
		intEnd, _ := strconv.Atoi(cidrEnd)
		// example: 192.168.1.1-1
		if endErr == nil && intEnd >= 1 && intEnd <= 255 {
			grep := regexp.MustCompile("(.*)\\.(.*)\\.(.*)\\.(.*)")
			cidrStart := grep.FindStringSubmatch(cidrIp)[4]
			intIpStart, _ := strconv.Atoi(cidrStart)
			intIpEnd, _ := strconv.Atoi(cidrEnd)
			if intIpStart == 0 || intIpEnd >= 255 {
				_, _, err := net.ParseCIDR(errorIp)
				if err != nil {
					return false
				}
			}
			return true
		} else if endErr != nil && intEnd == 0 { // example: 192.168.1.1-192.168.1.10
			grep := regexp.MustCompile("(.*)\\.(.*)\\.(.*)\\.(.*)")
			cidrStart := grep.FindStringSubmatch(cidrIp)[4]
			intIpStart, _ := strconv.Atoi(cidrStart)
			cidrEnd := grep.FindStringSubmatch(cidrEnd)[4]
			intIpEnd, _ := strconv.Atoi(cidrEnd)
			if intIpStart == 0 || intIpEnd >= 255 {
				_, _, err := net.ParseCIDR(errorIp)
				if err != nil {
					return false
				}
			}
			return true
		} else {
			_, _, err := net.ParseCIDR(errorIp)
			if err != nil {
				return false
			}
		}
	} else if len(regexp.MustCompile("/").FindAllStringIndex(target, -1)) == 1 {
		_, _, err := net.ParseCIDR(target)
		if err != nil {
			return false
		}
		return true
	} else {
		CustomHosts := strings.Split(target, ",")
		for _, checkHost := range CustomHosts {
			_, _, err := net.ParseCIDR(checkHost)
			if err != nil {
				return false
			}
		}
		return true
	}
	return true
}