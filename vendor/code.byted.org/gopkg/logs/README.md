# Golang日志库

[![Go Report Card](http://golang-report.byted.org/badge/code.byted.org/gopkg/logs)](http://golang-report.byted.org/report/code.byted.org/gopkg/logs)
[![build status](http://code.byted.org/ci/projects/10/status.png?ref=master)](http://code.byted.org/ci/projects/10?ref=master)


# 本日志库由于历史原因引入了较多的外部依赖，例如Metrics和databus，使用此日志库之前请周知

### 关于日志级别

比较有疑问的是以下几个日志级别的定位: Notice, Info, Trace。这里对日志级别的定义统一如下：

0. Trace
1. Debug
2. Info
3. Notice
4. Warn
5. Error
6. Fatal


### PushNotice
新版本中增加了对PushNotice的支持;

利用ctx缓存当前调用栈的kv信息, 因此需要提前调用NewNoticeCtx;

当栈返回时, 需要调用CtxFlushNotice;

下面是demo;

```go
func handler(ctx context.Context, req interface{}) (interface{}, error) {
	ctx = logs.NewNoticeCtx(ctx)
	defer logs.CtxFlushNotice(ctx)
	logs.CtxPushNotice(ctx, "method", "handler")
	return method0(ctx, req.(int))
}

func method0(ctx context.Context, id int) (interface{}, error) {
	logs.CtxPushNotice(ctx, "id", id)
	return method1(ctx)
}

func method1(ctx context.Context) (interface{}, error) {
	logs.CtxPushNotice(ctx, "method1", "this is method1")
	return nil, nil
}
```

### 新版升级: 增加metrics
新的logs增加了Warn, Error, Fatal三个字段的metrics.

新增的metrics为toutiao.service.log.{PSM}.throughput. 

metrics的tag中包含一个"level"字段, 用于显示错误等级, 分别为"WARNING", "ERROR", "CRITICAL", 为了和PY做统一, 所以稍有差异.


### 设计思想

日志模块分为logger和logger provider两个不同的组件, logger provider实现log往哪个地方写的逻辑。通常会有console, file, scribe等。

logger模块拥有自己的level，用于判断是否往各个provider输出日志，provider模块也有各自的level

Notice:
当前各个provider的实现在logs库里,这引入较多的外部依赖.为了更好地解耦这些依赖,我们将各个provider的实现放到了别的库里,如有需要新用户可以直接使用以下库provider,而不要使用logs库里的provider(FileProvider与ConsoleProvider除外)
- databus provider :code.byted.org/log_market/databus_provider.将日志写入databus.
- logAgent provider: code.byted.org/log_market/logagent_provider.将日志写入logAgent对接流式日志:https://docs.bytedance.net/doc/LeD6GaVSZgwpK1dVRLZj7f


#### 前缀

    level date time code

*注*
如果该日志库无法满足需求，，请尽量在自己使用的时候封装一层，而不要直接修改基础库

### optimize接口

在原版接口的基础上提供了高效的版的sf、fl、sffl接口，并且把原版只传format而无需格式化参数时的一些不必要开销剔除了。

sf接口：传入的格式化参数只支持string类型，从而避免反射消耗，占位语法为"%s"对应一个string参数，缺少的string参数不打印, 多余的string参数依次打印在末尾，%后跟其他字符会把这个字符打印，如果结尾为奇数个%不报错，最后的%被忽略，无其他复杂语法。
fl接口：支持第一个参数将日志的file location以string形式传入，从而避免location函数带来的巨大消耗。
sffl接口：包含以上两点特性。

### Processor定义
在log数据发送到各个provider之前，logger会把ctx, level, location, rawLog, kvs作为参数传递给provider
为了支持对logs在发送数据时的加密、特定转换函数，提供了一个processor接口，
用于实现在emit log的时候，fix raw log的展示，mask某些kvs的参数。
```go
/**
	Input:
		rawLog: 已经format过的log，即最终需要emit出来的函数
		kvs：contex中生成的和以KVs结尾的函数中的kvs列表
	Output:
		string: 处理过的log
		[]interface{}：处理过的kvs返回
		bool: 是否有效，如果为false，默认该条日志不处理了
 */
type Processor func(rawLog string, kvs ...interface{}) (string, []interface{}, bool)

```
其中，rawLog指的已经format过的log原文
kvs是最终emit的时候从context抽取和用户自定义的kv列表
