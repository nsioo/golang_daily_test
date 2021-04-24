package env

import (
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

const (
	UnknownIDC         = "-"
	DC_BR              = "br"
	DC_JP              = "jp"
	DC_IN              = "in"
	DC_FRAWS           = "fraws"  // 未确认的法兰克福机房
	DC_MAWSJP          = "mawsjp" // musical.ly 东京老机房
	DC_AGCQ            = "agcq"
	DC_AGGDSZ          = "aggdsz"
	DC_AGHBWH          = "aghbwh"
	DC_AGJSNJ          = "agjsnj"
	DC_AGLNSY          = "aglnsy"
	DC_AGSDQD          = "agsdqd"
	DC_AGSXLQ          = "agsxlq"
	DC_AGSXXA          = "agsxxa"
	DC_AGSY            = "agsy"
	DC_ALIEC1          = "aliec1"
	DC_ALIEC2          = "aliec2"
	DC_ALIGDSZ         = "aligdsz"
	DC_ALINC2          = "alinc2" // aliyun north
	DC_ALISC1          = "alisc1"
	DC_ALISG           = "alisg" // Singapore Aliyun
	DC_ALISH           = "alish"
	DC_ALIVA           = "aliva"
	DC_AWSBH           = "awsbh"
	DC_AWSBR           = "awsbr"
	DC_AWSFR           = "awsfr" // 法兰克福
	DC_AWSIN           = "awsin"
	DC_AWSJP           = "awsjp"
	DC_AWSJPGM         = "awsjpgm"
	DC_AWSNC1          = "awsnc1"
	DC_AWSNWC1         = "awsnwc1"
	DC_AWSSG           = "awssg"
	DC_AWSSGGM         = "awssggm"
	DC_AWSVAC          = "awsvac"
	DC_AWSVAGM         = "awsvagm"
	DC_BJGS            = "bjgs"
	DC_BJLGY           = "bjlgy"
	DC_BOE             = "boe"     // bytedance offline environment
	DC_IBOE            = "boei18n" // bytedance offline environment(international)
	DC_CA              = "ca"      // West America
	DC_COF             = "cof"
	DC_DEVBOX          = "devbox"
	DC_DEVBOXI18N      = "devboxi18n"
	DC_GALINC2         = "galinc2"
	DC_GALISG          = "galisg"
	DC_GALIVA          = "galiva"
	DC_GCPAU           = "gcpau"
	DC_GCPBE           = "gcpbe"
	DC_GCPBR           = "gcpbr"
	DC_GCPCA           = "gcpca"
	DC_GCPCH           = "gcpch"
	DC_GCPDE           = "gcpde"
	DC_GCPFI           = "gcpfi"
	DC_GCPGB           = "gcpgb"
	DC_GCPHK           = "gcphk"
	DC_GCPIN           = "gcpin"
	DC_GCPJPOSA        = "gcpjposa"
	DC_GCPJPTKY        = "gcpjptky"
	DC_GCPNL           = "gcpnl"
	DC_GCPSG           = "gcpsg"
	DC_GCPTW           = "gcptw"
	DC_GCPUSCBF        = "gcpuscbf"
	DC_GCPUSIAD        = "gcpusiad"
	DC_HKCJ            = "hkcj"
	DC_HL              = "hl"
	DC_HY              = "hy"
	DC_KSRU            = "ksru"
	DC_LF              = "lf"
	DC_LQ              = "lq" // 灵丘机房
	DC_MALIVA          = "maliva"
	DC_QDTOB           = "qdtob"
	DC_QDTOBIAAS       = "qdtobiaas"
	DC_SDQD            = "sdqd"
	DC_SG              = "sg"
	DC_SG1             = "sg1"
	DC_SGEA1           = "sgea1"
	DC_SGEA2           = "sgea2"
	DC_SGEE1           = "sgee1"
	DC_SGEE2           = "sgee2"
	DC_SGEE3           = "sgee3"
	DC_SGSAAS1LARKIDC1 = "sgsaas1larkidc1"
	DC_SGSAAS1LARKIDC2 = "sgsaas1larkidc2"
	DC_SGSAAS1LARKIDC3 = "sgsaas1larkidc3"
	DC_SYKA2LARK       = "syka2lark" // 顺义KA2机房(lark)
	DC_SYKA3LARK       = "syka3lark"
	DC_SYKALARK        = "sykalark" // 顺义KA机房(lark)
	DC_USEAST2A        = "useast2a" // 美东第二机房
	DC_USEAST3         = "useast3"
	DC_USWEST1A        = "uswest1a"
	DC_VA              = "va"
	DC_WJ              = "wj"
)

const (
	agentFileV4 = "/opt/tmp/consul_agent/netsegment"
	localFileV4 = "/opt/tiger/chadc/netsegment"
	agentFileV6 = "/opt/tmp/consul_agent/netsegment6"
	localFileV6 = "/opt/tiger/chadc/netsegment6"
)

const idcRefreshDur = regionRefreshDur

var (
	idc          atomic.Value // local idc (string)
	ipv4File     atomic.Value // file name (string)
	ipv6File     atomic.Value // file name (string)
	ipv4Segments atomic.Value // ipv4 segments - idc ([]*ipv4Segment)
	ipv6Segments atomic.Value // ipv6 segments - idc ([]*ipv6Segment)
)

// IDC .
func IDC() string {
	if v := idc.Load(); v != nil {
		return v.(string)
	}

	if dc := os.Getenv("RUNTIME_IDC_NAME"); dc != "" {
		idc.Store(dc)
		return dc
	}

	b, err := ioutil.ReadFile("/opt/tmp/consul_agent/datacenter")
	if err == nil {
		if dc := strings.TrimSpace(string(b)); dc != "" {
			idc.Store(dc)
			return dc
		}
	}

	cmd0 := exec.Command("/opt/tiger/consul_deploy/bin/determine_dc.sh")
	output0, err := cmd0.Output()
	if err == nil {
		dc := strings.TrimSpace(string(output0))
		if hasIDC(dc) {
			idc.Store(dc)
			return dc
		}
	}

	cmd := exec.Command(`bash`, `-c`, `sd report|grep "Data center"|awk '{print $3}'`)
	output, err := cmd.Output()
	if err == nil {
		dc := strings.TrimSpace(string(output))
		if hasIDC(dc) {
			idc.Store(dc)
			return dc
		}
	}

	idc.Store(UnknownIDC)
	return UnknownIDC
}

// GetIDCFromHost .
func GetIDCFromHost(ip string) string {
	netIP := net.ParseIP(ip)
	// ipv4
	if strings.Contains(ip, ".") {
		return lookupIpv4(netIP)
	}
	// ipv6
	if strings.Contains(ip, ":") {
		return lookupIpv6(netIP)
	}
	return UnknownIDC
}

// GetIDCList .
func GetIDCList() []string {
	return idcList()
}

// SetIDC .
func SetIDC(v string) {
	idc.Store(v)
}

// SetIDCFile only replace file ipv4 segment
func SetIDCFile(file string) {
	ipv4File.Store(file)
	updateIpv4Segments()
}

// SetIpv6File only replace file ipv6 segment6
func SetIpv6File(file string) {
	ipv6File.Store(file)
	updateIpv6Segments()
}

func init() {
	refreshIDCs()
}

// ATTENTION: IT COMES WITH A LOOP, DON'T CALL IT AGAIN.
func refreshIDCs() {
	defer time.AfterFunc(idcRefreshDur, refreshIDCs)
	updateIpv4Segments()
	updateIpv6Segments()
}

// updateIpv4Segments
func updateIpv4Segments() {
	var ipv4s []*ipv4Segment
	addItem := func(ipNet *net.IPNet, idc string) {
		ipv4s = append(ipv4s, newIpv4Segment(ipNet, idc))
	}

	fileV4s := []string{agentFileV4, localFileV4}
	readSegments(ipv4File, fileV4s, addItem)

	// sort & merge
	ipv4Segments.Store(mergeIpv4Segment(ipv4s))
}

// updateIpv6Segments
func updateIpv6Segments() {
	var ipv6s []*ipv6Segment
	addItem := func(ipNet *net.IPNet, idc string) {
		ipv6s = append(ipv6s, newIpv6Segment(ipNet, idc))
	}

	fileV6s := []string{agentFileV6, localFileV6}
	readSegments(ipv6File, fileV6s, addItem)

	// sort & merge
	ipv6Segments.Store(mergeIpv6Segment(ipv6s))
}

// readSegments add all read segments
func readSegments(userFile atomic.Value, defaultFiles []string, addItem func(ipNet *net.IPNet, idc string)) {
	var lines []string
	// read user's setting first
	file, _ := userFile.Load().(string)
	if file != "" {
		lines = readFile(file)
	}
	// if no user's setting, read default files
	if len(lines) == 0 {
		for _, file := range defaultFiles {
			lines = append(lines, readFile(file)...)
		}
	}
	// Deduplication: key=segment, value=idc
	var lineMap = make(map[string]string, len(lines))
	for _, line := range lines {
		s := strings.Split(line, " ")
		if len(s) != 2 {
			continue
		}
		// The former line has higher priority
		seg, idc := s[0], s[1]
		if _, ok := lineMap[seg]; ok {
			continue
		}
		_, ipNet, err := net.ParseCIDR(seg)
		if err != nil {
			continue
		}
		// add item
		addItem(ipNet, idc)
		lineMap[seg] = idc
	}
}

// readFile return lines
func readFile(file string) (lines []string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return lines
	}
	return strings.Split(string(content), "\n")
}

/* ipv4/ipv6 归属IDC 查询算法说明
 *
 * 本文件查询采用了 基于网段分片的二分查找算法
 * 我们知道 对于形如 10.10.8.0/22 对应的 ipv4 网段为 [10.10.8.0, 10.10.11.255], 换算为 uint32 为 [0xA0A0800, 0xA0A0BFF]
 * 因此, 对于给定的一组网段信息, 均可以转换为 uint32 的闭区间 形如 [a, b],[b+1,c],[d,e][e+1,f]....
 * 我们对上述闭区间从小到大排序, 则对于给定的 ipv4 地址, 则可以先转为 uint32, 在通过二分查找得到其网段归属, 获得该网段的 IDC 信息
 * ipv6 同理
 *
 * 算法效率
 * 基于二分查找的时间复杂度为 O(log(n)), 空间复杂度 O(n) 对于头条目前网段配置(ipv4 网段个数 60-70), 该算法最多查询 7 次
 * 附: 为什么不采用前缀树匹配 ?
 * 前缀树实际会很大, 如 10.10.8.0/22 有 22 位前缀, 如果再加上形如 10.224.0.0/18, 10.224.128.0/17, 10.11.0.0/16 的网段,
 * 则即使最优存储情况下, 前缀树需要保留 16-22 位(实际更多) 的分叉节点, 树的层级远大于 7, 并且存储结构比有序数组复杂, 因此不是最优选择
 */

// ipv4 section
//
// lookupIpv4 .
func lookupIpv4(ip net.IP) (idc string) {
	if len(ip) == net.IPv6len {
		ip = ip[12:16]
	}
	if len(ip) != net.IPv4len {
		return UnknownIDC
	}
	ipv4s, ok := ipv4Segments.Load().([]*ipv4Segment)
	if !ok || len(ipv4s) == 0 {
		return UnknownIDC
	}
	// Binary search
	target := bigEndianUint32(ip)
	search := func(i int) bool {
		return target <= ipv4s[i].end
	}
	index := sort.Search(len(ipv4s), search)
	// search not found
	if index < 0 || index >= len(ipv4s) {
		return UnknownIDC
	}
	if ipv4s[index].start <= target && target <= ipv4s[index].end {
		return ipv4s[index].idc
	} else {
		return UnknownIDC
	}
}

// ipv4Segment .
type ipv4Segment struct {
	start uint32
	end   uint32
	idc   string
}

// newIpv4Segment .
func newIpv4Segment(ipNet *net.IPNet, idc string) *ipv4Segment {
	segment := &ipv4Segment{
		idc:   idc,
		start: bigEndianUint32(ipNet.IP),
	}
	segment.end = segment.start + ^bigEndianUint32(ipNet.Mask)
	return segment
}

/* ipv4/ipv6 网段合并算法说明
 *
 * 由于网段配置文件采用了最长匹配原则, 形如 10.8.0.0/16, 10.8.0.0/17, 存在包含关系, 需要打平为互不包含的网段列表。
 * 我们先对原始网段配置, 按起始位置从小到大排序, 对于起始相同的, 先结束的小网段排在后;
 * 再遍历有序网段列表, 对于前后网段有包含关系的, 拆分成互不包含的网段, 由此得到打平后的有序网段配置。
 */

// mergeIpv4Segment Merge adjacent segments if idc is same
func mergeIpv4Segment(ipv4s []*ipv4Segment) (merged []*ipv4Segment) {
	if len(ipv4s) == 0 {
		return merged
	}
	sort.Sort(sortIpv4Segment(ipv4s))
	// 先打平所有网段
	var flatten = []*ipv4Segment{ipv4s[0]}
	var former, latter *ipv4Segment
	for formerIdx, latterIdx := 0, 1; latterIdx < len(ipv4s); {
		former, latter = flatten[formerIdx], ipv4s[latterIdx]
		// 判断前后两网段是否交叉, former.end >= latter.start
		if former.end >= latter.start {
			latterIdx++
			if former.start == latter.start {
				flatten[formerIdx] = latter
			} else {
				flatten[formerIdx] = &ipv4Segment{
					idc:   former.idc,
					start: former.start,
					end:   latter.start - 1,
				}
				flatten = append(flatten, latter)
			}
			// former.end 更后
			if former.end > latter.end {
				former.start = latter.end + 1
				flatten = append(flatten, former)
			}
		} else {
			// 两网段相互独立
			formerIdx++
			if formerIdx == len(flatten) {
				latterIdx++
				flatten = append(flatten, latter)
			}
		}
	}
	// 合并相邻网段
	former = flatten[0]
	merged = append(merged, flatten[0])
	for i := 1; i < len(flatten); i++ {
		latter = flatten[i]
		// idc 相同, 并且 former.end+1 == latter.start
		if former.idc == latter.idc && former.end+1 == latter.start {
			former.end = latter.end
		} else {
			former = latter
			merged = append(merged, latter)
		}
	}
	return merged
}

// sortIpv4Segment .
type sortIpv4Segment []*ipv4Segment

// Len .
func (s sortIpv4Segment) Len() int {
	return len(s)
}

// Swap .
func (s sortIpv4Segment) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less .
func (s sortIpv4Segment) Less(i, j int) bool {
	if s[i].start < s[j].start {
		return true
	}
	// 起始位置一样, 大网段在前
	if s[i].start == s[j].start && s[i].end > s[j].end {
		return true
	}
	return false
}

// ipv6 section
//
// lookupIpv6 .
func lookupIpv6(ip net.IP) (idc string) {
	if len(ip) != net.IPv6len {
		return UnknownIDC
	}
	ipv6s, ok := ipv6Segments.Load().([]*ipv6Segment)
	if !ok || len(ipv6s) == 0 {
		return UnknownIDC
	}
	// Binary search
	target := bigEndianUint128(ip)
	search := func(i int) bool {
		// target <= ipv6s[i].end
		return leU128(target, ipv6s[i].end)
	}
	index := sort.Search(len(ipv6s), search)
	// search not found
	if index < 0 || index >= len(ipv6s) {
		return UnknownIDC
	}
	// ipv6s[i].start <= target && target <= ipv6s[i].end
	if leU128(ipv6s[index].start, target) && leU128(target, ipv6s[index].end) {
		return ipv6s[index].idc
	} else {
		return UnknownIDC
	}
}

// ipv6 has 128 bits and uses [4]uint32.
type ipv6Segment struct {
	start uint128
	end   uint128
	idc   string
}

// newIpv6Segment .
func newIpv6Segment(ipNet *net.IPNet, idc string) *ipv6Segment {
	segment := &ipv6Segment{
		idc:   idc,
		start: bigEndianUint128(ipNet.IP),
	}
	segment.end = sumU128(segment.start, negU128(bigEndianUint128(ipNet.Mask)))
	return segment
}

// ipv6 地址遵循前缀匹配原则, 因此先基于网段起始位置排序, 在根据包含关系划分网段
func mergeIpv6Segment(ipv6s []*ipv6Segment) (merged []*ipv6Segment) {
	if len(ipv6s) == 0 {
		return merged
	}
	sort.Sort(sortIpv6Segment(ipv6s))
	// 打平所有网段
	var flatten = []*ipv6Segment{ipv6s[0]}
	var former, latter *ipv6Segment
	for formerIdx, latterIdx := 0, 1; latterIdx < len(ipv6s); {
		former, latter = flatten[formerIdx], ipv6s[latterIdx]
		// 判断前后两网段是否交叉, former.end >= latter.start
		if geU128(former.end, latter.start) {
			latterIdx++
			if eqU128(former.start, latter.start) {
				flatten[formerIdx] = latter
			} else {
				flatten[formerIdx] = &ipv6Segment{
					idc:   former.idc,
					start: former.start,
					end:   diffU128(latter.start, uint128{0, 0, 0, 1}),
				}
				flatten = append(flatten, latter)
			}
			// former.end 更后
			if gtU128(former.end, latter.end) {
				former.start = sumU128(latter.end, uint128{0, 0, 0, 1})
				flatten = append(flatten, former)
			}
		} else {
			// 两网段相互独立
			formerIdx++
			if formerIdx == len(flatten) {
				latterIdx++
				flatten = append(flatten, latter)
			}
		}
	}
	// 合并相邻网段
	former = flatten[0]
	merged = append(merged, flatten[0])
	for i := 1; i < len(flatten); i++ {
		latter = flatten[i]
		// idc 相同, 并且 former.end+1 == latter.start
		if former.idc == latter.idc && eqU128(sumU128(former.end, uint128{0, 0, 0, 1}), latter.start) {
			former.end = latter.end
		} else {
			former = latter
			merged = append(merged, latter)
		}
	}
	return merged
}

// sortIpv6Segment .
type sortIpv6Segment []*ipv6Segment

// Len .
func (s sortIpv6Segment) Len() int {
	return len(s)
}

// Swap .
func (s sortIpv6Segment) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less .
func (s sortIpv6Segment) Less(i, j int) bool {
	if ltU128(s[i].start, s[j].start) {
		return true
	}
	// 起始位置一样, 大网段在前
	if eqU128(s[i].start, s[j].start) && gtU128(s[i].end, s[j].end) {
		return true
	}
	return false
}

// uint128 Express with four uint32s
type uint128 [4]uint32

// geU128 return true when i >= j.
func geU128(i, j uint128) bool {
	return !ltU128(i, j)
}

// leU128 return true when i <= j.
func leU128(i, j uint128) bool {
	return !ltU128(j, i)
}

// eqU128 return true when i == j.
func eqU128(i, j uint128) bool {
	return i[0] == j[0] && i[1] == j[1] && i[2] == j[2] && i[3] == j[3]
}

// gtU128 return true when i > j.
func gtU128(i, j uint128) bool {
	return ltU128(j, i)
}

// ltU128 return true when i < j.
func ltU128(i, j uint128) bool {
	if i[0] != j[0] {
		return i[0] < j[0]
	}
	if i[1] != j[1] {
		return i[1] < j[1]
	}
	if i[2] != j[2] {
		return i[2] < j[2]
	}
	return i[3] < j[3]
}

// sumU128 .
func sumU128(i, j uint128) uint128 {
	var res uint128
	for k := 3; k > 0; k-- {
		tmp := uint64(i[k]) + uint64(j[k]) + uint64(res[k])
		res[k] = uint32(tmp)
		res[k-1] = uint32(tmp >> 32)
	}
	res[0] += i[0] + j[0]
	return res
}

// diffU128 .
func diffU128(i, j uint128) uint128 {
	var res uint128
	res[0], res[1], res[2], res[3] = i[0], i[1], i[2], i[3]
	for k := 3; k > 0; k-- {
		if res[k] >= j[k] {
			res[k] -= j[k]
		} else {
			tmp := uint128{0, 0, 0, 0}
			tmp[k-1] = 1
			res = diffU128(res, tmp)
			res[k] = uint32(uint64(1<<32) | uint64(res[k]) - uint64(j[k]))
		}
	}
	res[0] -= j[0]
	return res
}

// negU128 Negate number.
func negU128(i uint128) uint128 {
	return uint128{^i[0], ^i[1], ^i[2], ^i[3]}
}

// bigEndianUint32 .
func bigEndianUint32(p []byte) uint32 {
	_ = p[3]
	return uint32(p[0])<<24 | uint32(p[1])<<16 | uint32(p[2])<<8 | uint32(p[3])
}

// bigEndianUint128 .
func bigEndianUint128(p []byte) uint128 {
	_ = p[15]
	return uint128{bigEndianUint32(p[0:4]), bigEndianUint32(p[4:8]), bigEndianUint32(p[8:12]), bigEndianUint32(p[12:16])}
}
