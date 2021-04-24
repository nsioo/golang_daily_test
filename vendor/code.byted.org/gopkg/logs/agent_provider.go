package logs

import (
	"code.byted.org/log_market/gosdk"
	"os"
	"strconv"
)

const (
	canaryTaskName = "_canary"
	rpcTaskName    = "_rpc"
)

// AgentProvider : logagent provider, implement KVLogProvider
type AgentProvider struct {
	level    int
	isRPCLog bool
	pid      string
}

// 不再维护,建议使用https://code.byted.org/log_market/provider/ttlogagent.go里的NewLogAgentProvider()来替代
// NewAgentProvider : factory for AgentProvider.
func NewAgentProvider() *AgentProvider {
	return &AgentProvider{isRPCLog: false}
}

// NewRPCLogAgentProvider :
// 	if this provider is used to deal with RPC log, set isRPCLog true.
// 	only used for RPC log (in in kite/log.go).
func NewRPCLogAgentProvider() *AgentProvider {
	return &AgentProvider{isRPCLog: true}
}

// Init : implement KVLogProvider
func (ap *AgentProvider) Init() error {
	ap.pid = strconv.Itoa(os.Getpid())
	gosdk.Init()
	return nil
}

// SetLevel : implement KVLogProvider
func (ap *AgentProvider) SetLevel(level int) {
	ap.level = level
}

// WriteMsg : log agent do not support write message any more
func (ap *AgentProvider) WriteMsg(msg string, level int) error {
	// use WriteMsgKVs instead
	return nil
}

// WriteMsgKVs : implement KVLogProvider, core method for this provider
func (ap *AgentProvider) WriteMsgKVs(level int, msg string, headers map[string]string, kvs map[string]string) error {
	if level < ap.level {
		return nil
	}

	newKvs := make(map[string]string, len(kvs)+13)
	for k, v := range kvs {
		newKvs[k] = v
	}

	newKvs["_level"] = headers["level"]
	newKvs["_ts"] = headers["timestamp"]
	newKvs["_host"] = headers["host"]
	newKvs["_language"] = "go"
	newKvs["_psm"] = headers["psm"]
	newKvs["_cluster"] = headers["cluster"] // 对于RPC日志的最终KVs来说，cluster 是当前服务的集群，_cluster是远程服务的集群……不过这样设置真的OK么？
	newKvs["_logid"] = headers["logid"]
	newKvs["_deployStage"] = headers["stage"]
	newKvs["_podName"] = headers["pod_name"]
	newKvs["_process"] = ap.pid
	newKvs["_version"] = string(versionBytes) // "v1(6)"
	newKvs["_location"] = headers["location"]
	newKvs["_spanID"] = headers["span_id"]

	if ap.isRPCLog {
		newKvs["_taskName"] = rpcTaskName
	} else {
		newKvs["_taskName"] = headers["psm"]
	}

	message := &gosdk.Msg{
		Msg:  []byte(msg),
		Tags: newKvs,
	}

	gosdk.Send(newKvs["_taskName"], message)
	return nil
}

// Destroy : implement KVLogProvider
func (ap *AgentProvider) Destroy() error {
	gosdk.GracefullyExit()
	return nil
}

// Flush : implement KVLogProvider
func (ap *AgentProvider) Flush() error {
	return nil
}
