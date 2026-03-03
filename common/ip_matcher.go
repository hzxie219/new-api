package common

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// _trustedProxyCidrs holds the list of trusted reverse proxy IP ranges.
// Initialized once at startup via InitTrustedProxies; read-only at runtime.
var _trustedProxyCidrs []*net.IPNet

// InitTrustedProxies parses a comma-separated list of IP/CIDR strings and
// initializes the package-level _trustedProxyCidrs variable. Invalid entries
// are skipped with a WARN log. An empty string sets _trustedProxyCidrs to nil.
func InitTrustedProxies(proxies string) {
	if proxies == "" {
		_trustedProxyCidrs = nil
		return
	}

	parts := strings.Split(proxies, ",")
	result := make([]*net.IPNet, 0, len(parts))
	for _, part := range parts {
		s := strings.TrimSpace(part)
		if s == "" {
			continue
		}
		// Try parsing as CIDR first.
		_, network, err := net.ParseCIDR(s)
		if err == nil {
			result = append(result, network)
			continue
		}
		// Try as plain IP and auto-complete to /32 or /128.
		ip := net.ParseIP(s)
		if ip == nil {
			log.Printf("[WARN] InitTrustedProxies: invalid IP/CIDR %q, skipping", s)
			continue
		}
		bits := 32
		if ip.To4() == nil {
			bits = 128
		}
		cidr := fmt.Sprintf("%s/%d", s, bits)
		_, network, err = net.ParseCIDR(cidr)
		if err != nil {
			log.Printf("[WARN] InitTrustedProxies: failed to build CIDR for %q: %v, skipping", s, err)
			continue
		}
		result = append(result, network)
	}

	if len(result) == 0 {
		_trustedProxyCidrs = nil
	} else {
		_trustedProxyCidrs = result
	}
}

// ValidateCIDRList validates each entry in ips as a valid IP address or CIDR
// range. Returns the first error encountered, or nil for an empty list.
func ValidateCIDRList(ips []string) error {
	for _, s := range ips {
		_, _, err := net.ParseCIDR(s)
		if err == nil {
			continue
		}
		if net.ParseIP(s) != nil {
			continue
		}
		return fmt.Errorf("invalid IP/CIDR: %q", s)
	}
	return nil
}

// ParseCIDRList parses a slice of IP/CIDR strings into []*net.IPNet.
// Plain IP addresses are automatically completed to /32 (IPv4) or /128 (IPv6).
// Each call returns an independent slice with no shared state.
func ParseCIDRList(ips []string) ([]*net.IPNet, error) {
	result := make([]*net.IPNet, 0, len(ips))
	for _, s := range ips {
		_, network, err := net.ParseCIDR(s)
		if err == nil {
			result = append(result, network)
			continue
		}
		ip := net.ParseIP(s)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP/CIDR: %q", s)
		}
		bits := 32
		if ip.To4() == nil {
			bits = 128
		}
		cidr := fmt.Sprintf("%s/%d", s, bits)
		_, network, err = net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("failed to build CIDR for %q: %w", s, err)
		}
		result = append(result, network)
	}
	return result, nil
}

// IPMatchesCIDRList reports whether ipStr falls within any of the given CIDR
// ranges. Returns false for an invalid IP string or an empty cidrs list.
func IPMatchesCIDRList(ipStr string, cidrs []*net.IPNet) bool {
	if len(cidrs) == 0 {
		return false
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, cidr := range cidrs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// GetClientIP extracts the real client IP from a gin.Context.
// When the request's RemoteAddr belongs to a trusted proxy (_trustedProxyCidrs),
// it trusts the leftmost IP from the X-Forwarded-For header. Otherwise it uses
// RemoteAddr directly. Always returns a pure IP string without port.
func GetClientIP(c *gin.Context) string {
	remoteAddr := c.Request.RemoteAddr
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// RemoteAddr may already be a plain IP (no port).
		host = remoteAddr
	}

	remoteIP := net.ParseIP(host)
	if remoteIP != nil && isTrustedProxy(remoteIP) {
		xff := c.Request.Header.Get("X-Forwarded-For")
		if xff != "" {
			leftmost := strings.TrimSpace(strings.Split(xff, ",")[0])
			if leftmost != "" {
				return leftmost
			}
		}
	}

	return host
}

// isTrustedProxy reports whether ip is in the _trustedProxyCidrs list.
func isTrustedProxy(ip net.IP) bool {
	for _, cidr := range _trustedProxyCidrs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}
