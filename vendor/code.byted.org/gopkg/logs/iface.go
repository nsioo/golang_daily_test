package logs

import (
	"context"
)

// Logger接口对支持便捷组装的、不涉及反射的、显式传入fileLocation等功能的日志打印方法进行了抽象
// 这使得用户在使用此类功能时可以和具体的Logger类解耦，并使之统一化，将其收敛声明于本日志库中
// 为了兼容原有Logger接口，对其增加了新风格的别名

// logger的便捷组装日志的接口，提供基础的功能，支持interface格式化参数，所有的logger都应该实现该接口
type LoggerInterface interface {
	Fatal(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Notice(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Trace(format string, v ...interface{})
}

// 为统一风格额外声明，别名依赖go 1.9，暂不上线
type StdLogger interface {
	LoggerInterface
}

// logger 不涉及反射的sprintf接口，格式化参数只支持string，但组装会更加高效
type StdsfLogger interface {
	Fatalsf(format string, v ...string)
	Errorsf(format string, v ...string)
	Warnsf(format string, v ...string)
	Noticesf(format string, v ...string)
	Infosf(format string, v ...string)
	Debugsf(format string, v ...string)
	Tracesf(format string, v ...string)
}

// logger 显式传入fileLocation以消除location调用的显著消耗的接口
type StdflLogger interface {
	Fatalfl(fileLocation string, format string, v ...interface{})
	Errorfl(fileLocation string, format string, v ...interface{})
	Warnfl(fileLocation string, format string, v ...interface{})
	Noticefl(fileLocation string, format string, v ...interface{})
	Infofl(fileLocation string, format string, v ...interface{})
	Debugfl(fileLocation string, format string, v ...interface{})
	Tracefl(fileLocation string, format string, v ...interface{})
}

// logger 兼具sf和fl功能的接口
type StdsfflLogger interface {
	Fatalsffl(fileLocation string, format string, v ...string)
	Errorsffl(fileLocation string, format string, v ...string)
	Warnsffl(fileLocation string, format string, v ...string)
	Noticesffl(fileLocation string, format string, v ...string)
	Infosffl(fileLocation string, format string, v ...string)
	Debugsffl(fileLocation string, format string, v ...string)
	Tracesffl(fileLocation string, format string, v ...string)
}

// ctxLogger接口
type CtxLoggerInterface interface {
	CtxFatal(ctx context.Context, format string, v ...interface{})
	CtxError(ctx context.Context, format string, v ...interface{})
	CtxWarn(ctx context.Context, format string, v ...interface{})
	CtxNotice(ctx context.Context, format string, v ...interface{})
	CtxInfo(ctx context.Context, format string, v ...interface{})
	CtxDebug(ctx context.Context, format string, v ...interface{})
	CtxTrace(ctx context.Context, format string, v ...interface{})
}

// 为统一风格额外声明，别名依赖go 1.9，暂不上线
type CtxLogger interface {
	CtxLoggerInterface
}

// ctxLogger 不涉及反射的sprintf接口，格式化参数只支持string，但组装会更加高效
type CtxsfLogger interface {
	CtxFatalsf(ctx context.Context, format string, v ...string)
	CtxErrorsf(ctx context.Context, format string, v ...string)
	CtxWarnsf(ctx context.Context, format string, v ...string)
	CtxNoticesf(ctx context.Context, format string, v ...string)
	CtxInfosf(ctx context.Context, format string, v ...string)
	CtxDebugsf(ctx context.Context, format string, v ...string)
	CtxTracesf(ctx context.Context, format string, v ...string)
}

// ctxLogger 显式传入fileLocation以消除location调用的显著消耗的接口
type CtxflLogger interface {
	CtxFatalfl(fileLocation string, ctx context.Context, format string, v ...interface{})
	CtxErrorfl(fileLocation string, ctx context.Context, format string, v ...interface{})
	CtxWarnfl(fileLocation string, ctx context.Context, format string, v ...interface{})
	CtxNoticefl(fileLocation string, ctx context.Context, format string, v ...interface{})
	CtxInfofl(fileLocation string, ctx context.Context, format string, v ...interface{})
	CtxDebugfl(fileLocation string, ctx context.Context, format string, v ...interface{})
	CtxTracefl(fileLocation string, ctx context.Context, format string, v ...interface{})
}

// ctxLogger 兼具sf和fl功能的接口
type CtxsfflLogger interface {
	CtxFatalsffl(fileLocation string, ctx context.Context, format string, v ...string)
	CtxErrorsffl(fileLocation string, ctx context.Context, format string, v ...string)
	CtxWarnsffl(fileLocation string, ctx context.Context, format string, v ...string)
	CtxNoticesffl(fileLocation string, ctx context.Context, format string, v ...string)
	CtxInfosffl(fileLocation string, ctx context.Context, format string, v ...string)
	CtxDebugsffl(fileLocation string, ctx context.Context, format string, v ...string)
	CtxTracesffl(fileLocation string, ctx context.Context, format string, v ...string)
}

// ctxLoggerfmt接口，日志由KV对组装而成
type CtxKVsLogger interface {
	CtxFatalKVs(ctx context.Context, kvs ...interface{})
	CtxErrorKVs(ctx context.Context, kvs ...interface{})
	CtxWarnKVs(ctx context.Context, kvs ...interface{})
	CtxNoticeKVs(ctx context.Context, kvs ...interface{})
	CtxInfoKVs(ctx context.Context, kvs ...interface{})
	CtxDebugKVs(ctx context.Context, kvs ...interface{})
	CtxTraceKVs(ctx context.Context, kvs ...interface{})
}

// ctxLoggerfmt 显式传入fileLocation以消除location调用的显著消耗的接口
type CtxflKVsLogger interface {
	CtxFatalflKVs(fileLocation string, ctx context.Context, kvs ...interface{})
	CtxErrorflKVs(fileLocation string, ctx context.Context, kvs ...interface{})
	CtxWarnflKVs(fileLocation string, ctx context.Context, kvs ...interface{})
	CtxNoticeflKVs(fileLocation string, ctx context.Context, kvs ...interface{})
	CtxInfoflKVs(fileLocation string, ctx context.Context, kvs ...interface{})
	CtxDebugflKVs(fileLocation string, ctx context.Context, kvs ...interface{})
	CtxTraceflKVs(fileLocation string, ctx context.Context, kvs ...interface{})
}

// 旧ctxlogfmt接口，意义尚不明确
type CtxLogfmtInterface interface {
	CtxFatalKvs(ctx context.Context, kvs ...interface{})
	CtxErrorKvs(ctx context.Context, kvs ...interface{})
	CtxWarnKvs(ctx context.Context, kvs ...interface{})
	CtxNoticeKvs(ctx context.Context, kvs ...interface{})
	CtxInfoKvs(ctx context.Context, kvs ...interface{})
	CtxDebugKvs(ctx context.Context, kvs ...interface{})
	CtxTraceKvs(ctx context.Context, kvs ...interface{})
}
