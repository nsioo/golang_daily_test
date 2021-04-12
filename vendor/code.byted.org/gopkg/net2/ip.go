package net2

import (
	"bytes"
	"net"
)

var (
	v4InV6Prefix   = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}
	localIP        net.IP
	localIPStr     string
	localIPList    []net.IP
	localIPStrList []string
	privateNets    []net.IPNet
)

const UnknownIPAddr = "-"

func init() {
	for _, s := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "fc00::/7"} {
		_, n, _ := net.ParseCIDR(s)
		privateNets = append(privateNets, *n)
	}

	// get all network interfaces
	netIfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	// get all private ip list from non-loopback active net interfaces
	for _, netIface := range netIfaces {
		if netIface.Flags&net.FlagLoopback != 0 {
			// skip all Loopback interface
			continue
		}
		if netIface.Flags&net.FlagUp == 0 {
			// skip interface not UP
			continue
		}
		addrs, err := netIface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipnet.IP.IsLoopback() {
				continue
			}
			ip := ipnet.IP
			if IsPrivateIP(ip) {
				localIPList = append(localIPList, ip)
				localIPStrList = append(localIPStrList, ip.String())
			}
		}
	}

	// set local ip by priority of privateNets
	for _, pnet := range privateNets {
		for _, ip := range localIPList {
			if pnet.Contains(ip) {
				localIP = ip
				localIPStr = ip.String()
				return
			}
		}
	}
}

func GetLocalIP() net.IP {
	return localIP
}

func GetLocalIPStr() string {
	return localIPStr
}

// deprecated: use GetLocalIP or GetLocalIPStr
func GetLocalIp() string {
	if localIPStr == "" {
		return UnknownIPAddr
	}
	return localIPStr
}

func GetAllLocalIP() []net.IP {
	return localIPList
}

func GetAllLocalIPStr() []string {
	return localIPStrList
}

func IsPrivateIP(ip net.IP) bool {
	for _, n := range privateNets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func IsV4IP(ip net.IP) bool {
	if len(ip) == net.IPv4len || len(ip) == net.IPv6len && bytes.Equal(ip[:12], v4InV6Prefix) {
		return true
	}
	return false
}

func IsV6IP(ip net.IP) bool {
	if len(ip) == net.IPv6len && ! bytes.Equal(ip[:12], v4InV6Prefix) {
		return true
	}
	return false
}
