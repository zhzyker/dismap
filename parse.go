package dismap

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/malfunkt/iprange"
)

var DefaultPorts = []int{80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 443, 800, 801, 808, 880, 888, 889, 1000, 1080, 2601, 7001, 7007, 7010, 8000, 8001, 8002, 8003, 8004, 8005, 8006, 8007, 8008, 8009, 8010, 8011, 8012, 8016, 8017, 8018, 8019, 8022, 8029, 8030, 8060, 8069, 8070, 8080, 8081, 8082, 8083, 8084, 8085, 8086, 8087, 8088, 8089, 8090, 8091, 8092, 8093, 8094, 8095, 8096, 8097, 8098, 8099, 8100, 8101, 8108, 81110, 8161, 8175, 8188, 8189, 8200, 8222, 8300, 8360, 8443, 8445, 8448, 8484, 8800, 8848, 8879, 8880, 8881, 8888, 8899, 8983, 8989, 9000, 9001, 9002, 9008, 9010, 9043, 9060, 9080, 9081, 9082, 9083, 9084, 9085, 9086, 9087, 9088, 9089, 9090, 9091, 9092, 9093, 9094, 9095, 9096, 9097, 9098, 9099, 9100, 9200, 9443, 9448, 9800, 9981, 9986, 9988, 9998, 9999}

func ParseIPRange(Ips string) (ips []string, err error) {
	if Ips == "" {
		return ips, nil
	}
	defer func() {
		if e := recover(); e != nil {
			ips = nil
			err = fmt.Errorf("invalid ip-range : '%s'", Ips)
		}
	}()

	list, err := iprange.ParseList(Ips)
	if err != nil {
		return nil, err
	}
	rng := list.Expand()
	ips = make([]string, len(rng))
	for i := range rng {
		ips[i] = rng[i].String()
	}
	return ips, nil
}

func ParsePorts(Ports string) ([]int, error) {
	ports := make([]int, 0, 10)
	if Ports == "" {
		return DefaultPorts, nil
	}
	ranges := strings.Split(Ports, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid port selection segment: '%s'", r)
			}

			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid port number: '%s'", parts[0])
			}

			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid port number: '%s'", parts[1])
			}

			if p1 > p2 {
				return nil, fmt.Errorf("invalid port range: %d-%d", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}
		} else {
			port, err := strconv.Atoi(r)
			if err != nil {
				return nil, fmt.Errorf("invalid port number: '%s'", r)
			}
			ports = append(ports, port)
		}
	}
	return ports, nil
}

func ParseUrl(host string, port int) string {
	switch port {
	case 80:
		return "http://" + host
	case 443:
		return "https://" + host
	default:
		return fmt.Sprintf("http://%s:%d", host, port)
	}
}
