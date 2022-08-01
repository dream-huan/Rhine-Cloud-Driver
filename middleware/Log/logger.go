package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var pwd, _ = os.Getwd()
var baseURL = pwd + "/log/"
var logger *zap.SugaredLogger

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	file, err := os.OpenFile(baseURL+"access.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		file, err = os.Create(baseURL + "access.log")
		if err != nil {
			Errorw("创建日志文件失败")
		}
	}
	return zapcore.AddSync(file)
}

func init() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger = log.Sugar()
	defer logger.Sync()
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}
