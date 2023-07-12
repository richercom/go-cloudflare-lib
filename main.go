package main

import (
	"fmt"
	"net"
	"strings"
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
	ipStr := []string{"fe80::5448:a50a:f5d2:7892", "120.0.0.7"}
	fmt.Println(ConvertIP(ipStr))
}
