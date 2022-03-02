package parse

import (
	"github.com/zhzyker/dismap/configs"
	"regexp"
	"strconv"
	"strings"
)

func PortParse(port string) []int {
	var ports []int
	if port == "" {
		defPort := configs.DefaultPorts
		return defPort
	} else if len(regexp.MustCompile("-").FindAllStringIndex(port, -1)) == 1 {
		grep := regexp.MustCompile("(.*)-(.*)")
		portHead := grep.FindStringSubmatch(port)[1]
		portTail := grep.FindStringSubmatch(port)[2]
		intHead, _ := strconv.Atoi(portHead)
		intTail, _ := strconv.Atoi(portTail)
		for p := intHead; p <= intTail; p++ {
			ports = append(ports, p)

		}
		return ports
	} else {
		CustomPorts := strings.Split(port, ",")
		for _, port := range CustomPorts {
			p, _ := strconv.Atoi(port)
			ports = append(ports, p)
		}
		return ports
	}
}
