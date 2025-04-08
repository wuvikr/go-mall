package logger

import (
	"go-mall/common/enum"
	"go-mall/config"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _logger *zap.Logger

func init() {
	// 获取日志编码器
	encoder := getEncoder()
	// 获取文件日志写入器
	fileWriteSyncer := getFileLogWriter()

	// 根据环境选择日志输出方式
	var cores []zapcore.Core
	switch config.App.Env {
	case enum.ModeProd:
		// 生产环境只输出到文件中
		cores = append(cores, zapcore.NewCore(encoder, fileWriteSyncer, zapcore.InfoLevel))
	case enum.ModeDev:
		// 开发环境同时向控制台和文件输出日志， Debug级别的日志也会被输出
		cores = append(
			cores,
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
		)

	}
	core := zapcore.NewTee(cores...)
	_logger = zap.New(core)

}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	// encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
	// 使用 lumberjack 实现 logger rotate
	lumberJackLogger := &lumberjack.Logger{
		Filename:  config.App.Log.FilePath,
		MaxSize:   config.App.Log.FileMaxSize,      // 文件最大 100 M
		MaxAge:    config.App.Log.BackUpFileMaxAge, // 旧文件最多保留90天
		Compress:  false,
		LocalTime: true,
	}

	return zapcore.AddSync(lumberJackLogger)
}
