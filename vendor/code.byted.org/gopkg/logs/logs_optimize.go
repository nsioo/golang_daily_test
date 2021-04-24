package logs

import (
	"context"
)

// --------------------------------------------------------------------------------------------
// 	Optimized sf funcs without reflect
// --------------------------------------------------------------------------------------------

func Fatalsf(format string, v ...string) {
	defaultLogger.Fatalsf(format, v...)
}

func Errorsf(format string, v ...string) {
	defaultLogger.Errorsf(format, v...)
}

func Warnsf(format string, v ...string) {
	defaultLogger.Warnsf(format, v...)
}

func Noticesf(format string, v ...string) {
	defaultLogger.Noticesf(format, v...)
}

func Infosf(format string, v ...string) {
	defaultLogger.Infosf(format, v...)
}

func Debugsf(format string, v ...string) {
	if defaultLogger.level > LevelDebug {
		return
	}
	defaultLogger.Debugsf(format, v...)
}

func Tracesf(format string, v ...string) {
	defaultLogger.Tracesf(format, v...)
}

func CtxFatalsf(ctx context.Context, format string, v ...string) {
	defaultLogger.CtxFatalsf(ctx, format, v...)
}

func CtxErrorsf(ctx context.Context, format string, v ...string) {
	defaultLogger.CtxErrorsf(ctx, format, v...)
}

func CtxWarnsf(ctx context.Context, format string, v ...string) {
	defaultLogger.CtxWarnsf(ctx, format, v...)
}

func CtxNoticesf(ctx context.Context, format string, v ...string) {
	defaultLogger.CtxNoticesf(ctx, format, v...)
}

func CtxInfosf(ctx context.Context, format string, v ...string) {
	defaultLogger.CtxInfosf(ctx, format, v...)
}

func CtxDebugsf(ctx context.Context, format string, v ...string) {
	defaultLogger.CtxDebugsf(ctx, format, v...)
}

func CtxTracesf(ctx context.Context, format string, v ...string) {
	defaultLogger.CtxTracesf(ctx, format, v...)
}

// --------------------------------------------------------------------------------------------
// 	Funcs with locations param
// 	Used for go generate way, should not be called explicitly
// --------------------------------------------------------------------------------------------

func Fatalfl(fileLocation string, format string, v ...interface{}) {
	defaultLogger.Fatalfl(fileLocation, format, v...)
}

func Errorfl(fileLocation string, format string, v ...interface{}) {
	defaultLogger.Errorfl(fileLocation, format, v...)
}

func Warnfl(fileLocation string, format string, v ...interface{}) {
	defaultLogger.Warnfl(fileLocation, format, v...)
}

func Noticefl(fileLocation string, format string, v ...interface{}) {
	defaultLogger.Noticefl(fileLocation, format, v...)
}

func Infofl(fileLocation string, format string, v ...interface{}) {
	defaultLogger.Infofl(fileLocation, format, v...)
}

func Debugfl(fileLocation string, format string, v ...interface{}) {
	if defaultLogger.level > LevelDebug {
		return
	}
	defaultLogger.Debugfl(fileLocation, format, v...)
}

func Tracefl(fileLocation string, format string, v ...interface{}) {
	defaultLogger.Tracefl(fileLocation, format, v...)
}

func CtxFatalfl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxFatalfl(fileLocation, ctx, format, v...)
}

func CtxErrorfl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxErrorfl(fileLocation, ctx, format, v...)
}

func CtxWarnfl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxWarnfl(fileLocation, ctx, format, v...)
}

func CtxNoticefl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxNoticefl(fileLocation, ctx, format, v...)
}

func CtxInfofl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxInfofl(fileLocation, ctx, format, v...)
}

func CtxDebugfl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxDebugfl(fileLocation, ctx, format, v...)
}

func CtxTracefl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxTracefl(fileLocation, ctx, format, v...)
}

func CtxFatalflKvs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	defaultLogger.CtxFatalflKVs(fileLocation, ctx, kvs...)
}

func CtxErrorflKvs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	defaultLogger.CtxErrorflKVs(fileLocation, ctx, kvs...)
}

func CtxWarnflKvs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	defaultLogger.CtxWarnflKVs(fileLocation, ctx, kvs...)
}

func CtxNoticeflKvs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	defaultLogger.CtxNoticeflKVs(fileLocation, ctx, kvs...)
}

func CtxInfoflKvs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	defaultLogger.CtxInfoflKVs(fileLocation, ctx, kvs...)
}

func CtxDebugflKvs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	defaultLogger.CtxDebugflKVs(fileLocation, ctx, kvs...)
}

func CtxTraceflKvs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	defaultLogger.CtxTraceflKVs(fileLocation, ctx, kvs...)
}

func Fatalsffl(fileLocation string, format string, v ...string) {
	defaultLogger.Fatalsffl(fileLocation, format, v...)
}

func Errorsffl(fileLocation string, format string, v ...string) {
	defaultLogger.Errorsffl(fileLocation, format, v...)
}

func Warnsffl(fileLocation string, format string, v ...string) {
	defaultLogger.Warnsffl(fileLocation, format, v...)
}

func Noticesffl(fileLocation string, format string, v ...string) {
	defaultLogger.Noticesffl(fileLocation, format, v...)
}

func Infosffl(fileLocation string, format string, v ...string) {
	defaultLogger.Infosffl(fileLocation, format, v...)
}

func Debugsffl(fileLocation string, format string, v ...string) {
	if defaultLogger.level > LevelDebug {
		return
	}
	defaultLogger.Debugsffl(fileLocation, format, v...)
}

func Tracesffl(fileLocation string, format string, v ...string) {
	defaultLogger.Tracesffl(fileLocation, format, v...)
}

func CtxFatalsffl(fileLocation string, ctx context.Context, format string, v ...string) {
	defaultLogger.CtxFatalsffl(fileLocation, ctx, format, v...)
}

func CtxErrorsffl(fileLocation string, ctx context.Context, format string, v ...string) {
	defaultLogger.CtxErrorsffl(fileLocation, ctx, format, v...)
}

func CtxWarnsffl(fileLocation string, ctx context.Context, format string, v ...string) {
	defaultLogger.CtxWarnsffl(fileLocation, ctx, format, v...)
}

func CtxNoticesffl(fileLocation string, ctx context.Context, format string, v ...string) {
	defaultLogger.CtxNoticesffl(fileLocation, ctx, format, v...)
}

func CtxInfosffl(fileLocation string, ctx context.Context, format string, v ...string) {
	defaultLogger.CtxInfosffl(fileLocation, ctx, format, v...)
}

func CtxDebugsffl(fileLocation string, ctx context.Context, format string, v ...string) {
	defaultLogger.CtxDebugsffl(fileLocation, ctx, format, v...)
}

func CtxTracesffl(fileLocation string, ctx context.Context, format string, v ...string) {
	defaultLogger.CtxTracesffl(fileLocation, ctx, format, v...)
}
