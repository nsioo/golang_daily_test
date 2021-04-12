# Metrics

公司 metrics server client 的 Go 实现。

- 出于性能考虑数据异步，maxPendingSize=1000 or emitInterval=200ms 两个条件满足之一才发送；
    - 也可以使用 Flush 方法保证数据同步写到远程
- 在 `metrics.NewDefaultMetricsClientV2` 时指定 nocheck=true 可以忽略烦人的 DefineXXX 调用；
- Value 支持类型 float32 float64 int int8 int16 int32 int64 uint8 uint16 uint32 uint64 time.Duration
    - 其中 time.Duration 将表示为 nanosecond
- v1 与 v2 的区别 (MetricsClient vs MetricsClientV2):
    - v1 在 new 时要求namespace，而 emit 时除了 metrics name 还要输入 prefix，容易让人误用。v2 统一了在New时指定前缀，后面只能emit 后缀；
    - v1 tags 用的 map 结构，在高并发下，遍历map性能消耗高。v2 换成了 slice 而且如果没tags时可以省略；
    - 当前 v1 底层实际用了 v2 的逻辑，tags map 转 slice 时为了内部优化做了sort，有额外消耗;
    - 新的项目应该都使用 v2;
    - 对于v2: 请保证没有重复tag name，否则metrics查询出来的结果是未定义的；
- 性能优化提示
    - 对QPS高，又是动态tag组合的场景无解, 请确保tag组合是可以预定义的;
    - 通过 RegisterCounter / RegisterTimer 的方式预定义metrics，可以大量减少每次Emit导致的额外CPU计算消耗；
    - 如不通过预定义metrics，请尽量在业务侧对数据进行累加(如通过int64定期reset和emit);
    - 如不通过预定义metrics, 对于Timer类的数据，当前metrics server要求需要全部数据都发往server，而协议比较低效;
    - 重复: 请使用 MetricsClientV2, 而不是 MetricsClient, MetricsClient 主要为了兼容历史代码;
- 示例代码，详见 example/main.go
- [Metrics系统使用说明](https://bytedance.feishu.cn/docs/GHFmzle2R6a7cGvAqMWlbc#)

# Metrics V3

公司 metrics server client 的下一代 Go 实现。解决了 V2 中的丢点问题，并发性能较 V2 有明显的提升，同时还保留了一些动态能力，同时保证预聚合后的打点不会丢弃。V3 与 V2，V1 的 API 不兼容，可以与 V1 V2 一同使用。

```
 ↳ go test -test.v -test.bench ^BenchmarkClientParallel$ -test.run ^$
goos: darwin
goarch: amd64
pkg: code.byted.org/gopkg/metrics/v3
BenchmarkClientParallel
BenchmarkClientParallel/random_access
BenchmarkClientParallel/random_access-12         	 2815730	       419 ns/op	      89 B/op	       1 allocs/op
BenchmarkClientParallel/hot_key_hit
BenchmarkClientParallel/hot_key_hit-12           	 9189608	       134 ns/op

 ↳ go test -test.v -test.bench ^BenchmarkClientV2Parallel$ -test.run ^$
goos: darwin
goarch: amd64
pkg: code.byted.org/gopkg/metrics
BenchmarkClientV2Parallel
BenchmarkClientV2Parallel/random_access
BenchmarkClientV2Parallel/random_access-12         	  317238	      3831 ns/op	     307 B/op	       1 allocs/op
BenchmarkClientV2Parallel/hot_key_hit
BenchmarkClientV2Parallel/hot_key_hit-12           	  621397	      1970 ns/op	       0 B/op	       0 allocs/op

```

## Feature

- 支持动态的 Tag value；
- 支持 SDK 侧多值，一次打点可以同时打多个不同或者相同类型的指标，减少重复开销；
- 缓存并预聚合 counter store 类型的打点，Metrics 序列缓存从 1000 提升到了 1048575((1<<20)-1) 个，到达最大值前缓存不会被清理，提升了缓存效果；
- 异步发送打点，触发清理时或者每 500 毫秒异步发送全部缓存序列；
- 不再丢弃预聚合的 counter 打点；

### 相比 V2 Client 的区别

- V2 使用 `sync.Map` 缓存序列，V3 使用基于原子操作的无锁结构提升并发效率；
- 仅丢弃未聚合的 timer 类型打点，不丢 counter 及其它类型，保证打点准确性；
- 预聚合更多类型打点并且缓存周期更长，因此聚合效果更好；
- tag values 支持动态声明，移除 V2 最佳性能实践时要求静态 value 的约束；

#### 缺点
- V3 要求预对齐（预声明）Tag keys;

## Install

`go get -u code.byted.org/gopkg/metrics/v3`

## Example

```go
package example

import (
	"fmt"

	m "code.byted.org/gopkg/metrics/v3"
)

func Example() {
	// Initialize client with options.
	client := m.NewClient(
		"metrics.sdk",
		// Client options.
		m.SetTceTags(), m.SetGlobalTags(m.T{Name: "hello", Value: "world"}),
	)
	// Close as gracefully exit.
	defer client.Close()

	// Declare a new metric with tag keys, tag values is not declaration required.
	metric := client.NewMetric("test", []string{"foo", "bar", "baz"}...)

	// Send both timer and counter metrics with the same tag keys as above, tag value can be grabbed in runtime.
	tags := []m.T{m.T{Name: "foo", Value: "a"}, m.T{Name: "bar", Value: "b"}, m.T{Name: "baz", Value: "c"}}
	err := metric.WithTags(tags...).Emit(
		// A rate counter metric with the default suffix "rate"
		m.Incr(1),
		// Another rate counter type metric with the suffix "send-size" at the tail of the metric name
		m.WithSuffix("send_size").Incr(1),
		// A timer type metric with the suffix "latency" at the tail of the metric name
		m.WithSuffix("latency").Observe(100),
		// It is OK to use multiple counter metrics. 
		m.WithSuffix("recv_size").Incr(1),
	)
	if err != nil {
		fmt.Printf("Emit metrics error: %s", err.Error())
	}
}
```

## Document

GoDoc: https://codebase.byted.org/godoc/code.byted.org/gopkg/metrics/v3/

Design Document: https://bytedance.feishu.cn/docs/doccnmZIWzCpNrJPQdrkGlIhm2e#
