package main

import (
	"errors"
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

func ConvertIP(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", errors.New("")
	}
	if ip.To4() != nil {
		return ip.String(), nil
	}
	// v6 to cidr
	cidr := convertIPv6ToCIDR(ip.String(), 64)
	return cidr, nil

}

func main() {
	ipStr := "fe80::5448:a50a:f5d2:7892"
	fmt.Println(ConvertIP(ipStr))
}
