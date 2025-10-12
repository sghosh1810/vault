package validation

import (
	"fmt"
	"net"
	"strings"
)

func IsIPInRange(ipStr, cidrStr string) any {
	ipStr = normalizeIP(ipStr)
	ip := net.ParseIP(ipStr)

	fmt.Println(ipStr, cidrStr)
	if ip == nil {
		return map[string]string{
			"status":  "fail",
			"message": "Access valid for internal ips only",
		}
	}

	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return map[string]string{
			"status":  "fail",
			"message": "Access valid for internal ips only",
		}
	}

	if ipNet.Contains(ip) {
		return nil
	} else {
		return map[string]string{
			"status":  "fail",
			"message": "Access valid for internal ips only",
		}
	}
}

func normalizeIP(ip string) string {
	if ip == "::1" {
		return "127.0.0.1"
	}

	// Handle IPv6-mapped IPv4 addresses (e.g., ::ffff:10.0.0.1)
	if after, ok := strings.CutPrefix(ip, "::ffff:"); ok {
		return after
	}

	// Normalize zone suffix (e.g., fe80::1%eth0)
	if i := strings.IndexByte(ip, '%'); i != -1 {
		ip = ip[:i]
	}

	// Try to parse to ensure it’s valid
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ""
	}

	return ip
}
