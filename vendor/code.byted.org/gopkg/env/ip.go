package env

import (
	"net"
	"os"

	"code.byted.org/gopkg/net2"
)

var inTCE bool
var (
	tceIPV4  net.IP
	tceIPV6  net.IP
	hostIPV4 net.IP
	hostIPV6 net.IP

	tceIPV4Str  string
	tceIPV6Str  string
	hostIPV4Str string
	hostIPV6Str string
)

func init() {
	if os.Getenv("IS_TCE_DOCKER_ENV") == "1" {
		inTCE = true
		tceIPV4Str = os.Getenv("MY_HOST_IP")
		if tceIPV4Str != "" {
			tceIPV4 = net.ParseIP(tceIPV4Str)
		}
		tceIPV6Str = os.Getenv("MY_HOST_IPV6")
		if tceIPV6Str != "" {
			tceIPV6 = net.ParseIP(tceIPV6Str)
		}
	}

	ips := net2.GetAllLocalIP()
	for i := range ips {
		if hostIPV4 == nil && net2.IsV4IP(ips[i]) {
			hostIPV4 = ips[i]
			hostIPV4Str = hostIPV4.String()
		}
		if hostIPV6 == nil && net2.IsV6IP(ips[i]) {
			hostIPV6 = ips[i]
			hostIPV6Str = hostIPV6.String()
		}
	}
}

// HostIP 保留接口，为了兼容性优先返回 IPV4，其次是 IPV6
func HostIP() string {
	if HostIPV4() != "" {
		return HostIPV4()
	}
	return HostIPV6()
}

// HostIPV4 返回 ipv4 的 string
// 如果没有 ipv4，那么会返回空字符串""
// 如果有多个 ipv4，会随机返回其中一个
// 如果在多个 ip 情况之下想要获取所有的 ip，可以使用 net2.GetAllLocalIP
func HostIPV4() string {
	if inTCE && tceIPV4Str != "" {
		return tceIPV4Str
	}

	return hostIPV4Str
}

// HostIPV6 返回 ipv6 的 string
// 表现同上
func HostIPV6() string {
	if inTCE && tceIPV6Str != "" {
		return tceIPV6Str
	}

	return hostIPV6Str
}

// IPV4 返回 ipv4
// 表现同上
func IPV4() net.IP {
	if inTCE && tceIPV4 != nil {
		return tceIPV4
	}

	return hostIPV4
}

// IPV6 返回 ipv6
// 表现同上
func IPV6() net.IP {
	if inTCE && tceIPV6 != nil {
		return tceIPV6
	}

	return hostIPV6
}

// IP 优先返回 IPV6，如果没有 IPV6 则返回 IPV4
// 适用于 不在意获取到的是 IPV4 还是 IPV6 的业务
func IP() net.IP {
	if IPV6() != nil {
		return IPV6()
	}
	return IPV4()
}

// IPStr 同上，返回 string
func IPStr() string {
	if HostIPV6() != "" {
		return HostIPV6()
	}
	return HostIPV4()
}

func HasIPV4() bool {
	return (inTCE && tceIPV4 != nil) || hostIPV4 != nil
}

func HasIPV6() bool {
	return (inTCE && tceIPV6 != nil) || hostIPV6 != nil
}

func IsDualStack() bool {
	return HasIPV4() && HasIPV6()
}

func IsIPV4Only() bool {
	return HasIPV4() && !HasIPV6()
}

func IsIPV6Only() bool {
	return !HasIPV4() && HasIPV6()
}
