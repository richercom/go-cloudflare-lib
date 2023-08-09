package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

func convertIPv6ToCIDRF(ipv6CIDR string) string {
	ip, ipNet, err := net.ParseCIDR(ipv6CIDR)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(ip, ipNet)
	ipStr := ip.String()
	prefixLength, _ := ipNet.Mask.Size()
	prefix := strings.Split(ipStr, ":")[:prefixLength/16]
	cidrAddress := fmt.Sprintf("%s:/64", strings.Join(prefix, ":"))
	return cidrAddress
}

func convertIPv6ToCIDR(ipv6Address string, prefixLen int) string {
	ip := net.ParseIP(ipv6Address)
	ip = ip.Mask(net.CIDRMask(prefixLen, 128))
	return fmt.Sprintf("%s/%d", ip.String(), prefixLen)
}

func ConvertIP(ipStr []string) (newIPs []string) {
	if len(ipStr) < 1 {
		return
	}
	for i, _ := range ipStr {
		ip := net.ParseIP(ipStr[i])
		if ip == nil {
			fmt.Println("format ip err.", ip)
			continue
		}
		if ip.To4() != nil {
			newIPs = append(newIPs, ip.String())
			continue
		}
		// v6 to cidr
		cidr := convertIPv6ToCIDR(ip.String(), 64)
		newIPs = append(newIPs, cidr)
	}
	return
}

func main() {

	domain := "www.example.com"
	re := regexp.MustCompile(`[a-zA-Z0-9-]+\.([a-zA-Z]+)$`)
	match := re.FindStringSubmatch(domain)
	if len(match) > 1 {
		suffix := match[0]
		fmt.Println("Domain suffix:", suffix)
	} else {
		fmt.Println("Invalid domain format")
	}

	// 假设有两个时间戳 t1 和 t2
	t1 := time.Unix(1630100000, 0) // 第一个时间戳
	t2 := time.Unix(1630103000, 0) // 第二个时间戳
	fmt.Println(t1, t2)

	// 定义5分钟的时间间隔

	fiveMinutes := time.Duration(5) * time.Minute

	// 计算两个时间戳的差值
	diff := t2.Sub(t1)
	fmt.Println(diff)
	fmt.Println("555:", fiveMinutes)

	// 判断差值是否等于5分钟
	if diff > fiveMinutes {
		fmt.Println("两个时间戳相差 大于 5分钟")
	} else {
		fmt.Println("两个时间戳不超过5分钟")
	}

	ipStr := []string{"fe80::5448:a50a:f5d2:7892", "120.0.0.7"}
	fmt.Println(ConvertIP(ipStr))
}
