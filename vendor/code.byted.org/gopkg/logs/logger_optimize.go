package logs

import (
	"context"
	"fmt"
	"sync/atomic"

	"code.byted.org/gopkg/logs/utils"
)

// --------------------------------------------------------------------------------------------
// 	Optimized sf funcs without reflect
// --------------------------------------------------------------------------------------------

func (l *Logger) Fatalsf(format string, v ...string) {
	if LevelFatal < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelFatal, "", format)
		return
	}
	l.fmtLog(context.Background(), LevelFatal, "", utils.Sprintf(format, v...))
}

func (l *Logger) Errorsf(format string, v ...string) {
	if LevelError < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelError, "", format)
		return
	}
	l.fmtLog(context.Background(), LevelError, "", utils.Sprintf(format, v...))
}

func (l *Logger) Warnsf(format string, v ...string) {
	if LevelWarn < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelWarn, "", format)
		return
	}
	l.fmtLog(context.Background(), LevelWarn, "", utils.Sprintf(format, v...))
}

func (l *Logger) Noticesf(format string, v ...string) {
	if LevelNotice < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelNotice, "", format)
		return
	}
	l.fmtLog(context.Background(), LevelNotice, "", utils.Sprintf(format, v...))
}

func (l *Logger) Infosf(format string, v ...string) {
	if LevelInfo < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelInfo, "", format)
		return
	}
	l.fmtLog(context.Background(), LevelInfo, "", utils.Sprintf(format, v...))
}

func (l *Logger) Debugsf(format string, v ...string) {
	if LevelDebug < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelDebug, "", format)
		return
	}
	l.fmtLog(context.Background(), LevelDebug, "", utils.Sprintf(format, v...))
}

func (l *Logger) Tracesf(format string, v ...string) {
	if LevelTrace < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelTrace, "", format)
		return
	}
	l.fmtLog(context.Background(), LevelTrace, "", utils.Sprintf(format, v...))
}

func (l *Logger) CtxFatalsf(ctx context.Context, format string, v ...string) {
	if LevelFatal < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelFatal, "", format)
		return
	}
	l.fmtLog(ctx, LevelFatal, "", utils.Sprintf(format, v...))
}

func (l *Logger) CtxErrorsf(ctx context.Context, format string, v ...string) {
	if LevelError < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelError, "", format)
		return
	}
	l.fmtLog(ctx, LevelError, "", utils.Sprintf(format, v...))
}

func (l *Logger) CtxWarnsf(ctx context.Context, format string, v ...string) {
	if LevelWarn < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelWarn, "", format)
		return
	}
	l.fmtLog(ctx, LevelWarn, "", utils.Sprintf(format, v...))
}

func (l *Logger) CtxNoticesf(ctx context.Context, format string, v ...string) {
	if LevelNotice < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelNotice, "", format)
		return
	}
	l.fmtLog(ctx, LevelNotice, "", utils.Sprintf(format, v...))
}

func (l *Logger) CtxInfosf(ctx context.Context, format string, v ...string) {
	if LevelInfo < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelInfo, "", format)
		return
	}
	l.fmtLog(ctx, LevelInfo, "", utils.Sprintf(format, v...))
}

func (l *Logger) CtxDebugsf(ctx context.Context, format string, v ...string) {
	if LevelDebug < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelDebug, "", format)
		return
	}
	l.fmtLog(ctx, LevelDebug, "", utils.Sprintf(format, v...))
}

func (l *Logger) CtxTracesf(ctx context.Context, format string, v ...string) {
	if LevelTrace < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelTrace, "", format)
		return
	}
	l.fmtLog(ctx, LevelTrace, "", utils.Sprintf(format, v...))
}

// --------------------------------------------------------------------------------------------
// 	Funcs with locations param
// 	Used for go generate way, should not be called explicitly
// --------------------------------------------------------------------------------------------

//----------------------------funcs by reflect----------------------------

func (l *Logger) Fatalfl(fileLocation string, format string, v ...interface{}) {
	if LevelFatal < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelFatal, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelFatal, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) Errorfl(fileLocation string, format string, v ...interface{}) {
	if LevelError < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelError, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelError, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) Warnfl(fileLocation string, format string, v ...interface{}) {
	if LevelWarn < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelWarn, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelWarn, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) Noticefl(fileLocation string, format string, v ...interface{}) {
	if LevelNotice < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelNotice, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelNotice, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) Infofl(fileLocation string, format string, v ...interface{}) {
	if LevelInfo < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelInfo, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelInfo, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) Debugfl(fileLocation string, format string, v ...interface{}) {
	if LevelDebug < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelDebug, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelDebug, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) Tracefl(fileLocation string, format string, v ...interface{}) {
	if LevelTrace < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelTrace, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelTrace, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) CtxFatalfl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	if LevelFatal < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelFatal, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelFatal, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) CtxErrorfl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	if LevelError < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelError, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelError, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) CtxWarnfl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	if LevelWarn < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelWarn, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelWarn, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) CtxNoticefl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	if LevelNotice < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelNotice, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelNotice, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) CtxInfofl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	if LevelInfo < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelInfo, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelInfo, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) CtxDebugfl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	if LevelDebug < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelDebug, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelDebug, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) CtxTracefl(fileLocation string, ctx context.Context, format string, v ...interface{}) {
	if LevelTrace < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelTrace, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelTrace, fileLocation, fmt.Sprintf(format, v...))
}

func (l *Logger) CtxFatalflKVs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	if LevelFatal < atomic.LoadInt32(&l.level) {
		return
	}
	l.fmtLog(ctx, LevelFatal, fileLocation, "", kvs...)
}

func (l *Logger) CtxErrorflKVs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	if LevelError < atomic.LoadInt32(&l.level) {
		return
	}
	l.fmtLog(ctx, LevelError, fileLocation, "", kvs...)
}

func (l *Logger) CtxWarnflKVs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	if LevelWarn < atomic.LoadInt32(&l.level) {
		return
	}
	l.fmtLog(ctx, LevelWarn, fileLocation, "", kvs...)
}

func (l *Logger) CtxNoticeflKVs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	if LevelNotice < atomic.LoadInt32(&l.level) {
		return
	}
	l.fmtLog(ctx, LevelNotice, fileLocation, "", kvs...)
}

func (l *Logger) CtxInfoflKVs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	if LevelInfo < atomic.LoadInt32(&l.level) {
		return
	}
	l.fmtLog(ctx, LevelInfo, fileLocation, "", kvs...)
}

func (l *Logger) CtxDebugflKVs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	if LevelDebug < atomic.LoadInt32(&l.level) {
		return
	}
	l.fmtLog(ctx, LevelDebug, fileLocation, "", kvs...)
}

func (l *Logger) CtxTraceflKVs(fileLocation string, ctx context.Context, kvs ...interface{}) {
	if LevelTrace < atomic.LoadInt32(&l.level) {
		return
	}
	l.fmtLog(ctx, LevelTrace, fileLocation, "", kvs...)
}

//----------------------------funcs without reflect----------------------------

func (l *Logger) Fatalsffl(fileLocation string, format string, v ...string) {
	if LevelFatal < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelFatal, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelFatal, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) Errorsffl(fileLocation string, format string, v ...string) {
	if LevelError < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelError, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelError, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) Warnsffl(fileLocation string, format string, v ...string) {
	if LevelWarn < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelWarn, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelWarn, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) Noticesffl(fileLocation string, format string, v ...string) {
	if LevelNotice < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelNotice, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelNotice, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) Infosffl(fileLocation string, format string, v ...string) {
	if LevelInfo < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelInfo, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelInfo, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) Debugsffl(fileLocation string, format string, v ...string) {
	if LevelDebug < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelDebug, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelDebug, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) Tracesffl(fileLocation string, format string, v ...string) {
	if LevelTrace < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(context.Background(), LevelTrace, fileLocation, format)
		return
	}
	l.fmtLog(context.Background(), LevelTrace, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) CtxFatalsffl(fileLocation string, ctx context.Context, format string, v ...string) {
	if LevelFatal < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelFatal, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelFatal, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) CtxErrorsffl(fileLocation string, ctx context.Context, format string, v ...string) {
	if LevelError < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelError, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelError, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) CtxWarnsffl(fileLocation string, ctx context.Context, format string, v ...string) {
	if LevelWarn < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelWarn, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelWarn, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) CtxNoticesffl(fileLocation string, ctx context.Context, format string, v ...string) {
	if LevelNotice < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelNotice, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelNotice, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) CtxInfosffl(fileLocation string, ctx context.Context, format string, v ...string) {
	if LevelInfo < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelInfo, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelInfo, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) CtxDebugsffl(fileLocation string, ctx context.Context, format string, v ...string) {
	if LevelDebug < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelDebug, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelDebug, fileLocation, utils.Sprintf(format, v...))
}

func (l *Logger) CtxTracesffl(fileLocation string, ctx context.Context, format string, v ...string) {
	if LevelTrace < atomic.LoadInt32(&l.level) {
		return
	}
	if len(v) == 0 {
		l.fmtLog(ctx, LevelTrace, fileLocation, format)
		return
	}
	l.fmtLog(ctx, LevelTrace, fileLocation, utils.Sprintf(format, v...))
}
