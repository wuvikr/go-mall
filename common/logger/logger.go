package logger

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	ctx     context.Context
	traceId string
	spanId  string
	pSpanId string
	_logger *zap.Logger
}

func New(ctx context.Context) *logger {
	var traceId, spanId, pSpanId string
	if ctx.Value("traceid") != nil {
		traceId = ctx.Value("traceid").(string)
	}
	if ctx.Value("spanid") != nil {
		spanId = ctx.Value("spanid").(string)
	}
	if ctx.Value("pspanid") != nil {
		pSpanId = ctx.Value("pspanid").(string)
	}

	return &logger{
		ctx:     ctx,
		traceId: traceId,
		spanId:  spanId,
		pSpanId: pSpanId,
		_logger: _logger,
	}
}

func (l *logger) log(lvl zapcore.Level, msg string, kv ...any) {

	// 如果kv的长度不是偶数，则添加一个unknown
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}

	// 日志中添加追踪id
	kv = append(kv, "traceid", l.traceId, "spanid", l.spanId, "pspanid", l.pSpanId)

	// 增加日志调用者信息, 方便查日志时定位程序位置
	funcName, file, line := l.getLoggerCaller()
	kv = append(kv, "func", funcName, "file", file, "line", line)

	// 将kv转换为zap.Field
	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		key := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(key, kv[i+1]))
	}

	ce := l._logger.Check(lvl, msg)
	ce.Write(fields...)
}

func (l *logger) Debug(msg string, kv ...any) {
	l.log(zapcore.DebugLevel, msg, kv...)
}

func (l *logger) Info(msg string, kv ...any) {
	l.log(zapcore.InfoLevel, msg, kv...)
}

func (l *logger) Warn(msg string, kv ...any) {
	l.log(zapcore.WarnLevel, msg, kv...)
}

func (l *logger) Error(msg string, kv ...any) {
	l.log(zapcore.ErrorLevel, msg, kv...)
}

func (l *logger) getLoggerCaller() (funcName, file string, line int) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return
	}

	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}
