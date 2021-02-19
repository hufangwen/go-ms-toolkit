/*
 @Desc

 @Date 2021-01-27 17:20
 @Author inori
*/
package log

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	KeyTraceId  = "traceId"
	KeyParentId = "parentId"
	KeySpanId   = "spanId"
)

type RotateConfig struct {
	FilePath   string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

type Config struct {
	Rotate   RotateConfig
	FilePath string
	FileName string
	Level    string
}

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger

	_logLevelMap = map[string]zapcore.Level{
		"panic": zapcore.PanicLevel,
		"fatal": zapcore.FatalLevel,
		"error": zapcore.ErrorLevel,
		"warn":  zapcore.WarnLevel,
		"info":  zapcore.InfoLevel,
		"debug": zapcore.DebugLevel,
	}
	config = Config{
		Rotate: RotateConfig{
			FilePath:   "/data/logs",
			Filename:   "default.log",
			MaxSize:    200,
			MaxBackups: 7,
			MaxAge:     7,
		},
		Level: "debug",
	}
)

func Init(conf ...Config) {
	if conf != nil && len(conf) > 0 {
		config = conf[0]
	}
	rotateConfig := config.Rotate
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	instanceID := os.Getenv("INSTANCE_ID")

	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(rotateConfig.FilePath, instanceID, rotateConfig.Filename),
		MaxSize:    rotateConfig.MaxSize,
		MaxBackups: rotateConfig.MaxBackups,
		MaxAge:     rotateConfig.MaxAge,
		Compress:   false,
	}
	writeSync := zapcore.AddSync(lumberJackLogger)
	core := zapcore.NewCore(encoder, writeSync, _logLevelMap[config.Level])
	logger = zap.New(core)
	sugar = logger.Sugar()
	defer func() { _ = logger.Sync() }() // flushes buffer, if any
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	sugar.Debugw(fmt.Sprintf(format, args...), buildTraceKeyValues(ctx)...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	sugar.Infow(fmt.Sprintf(format, args...), buildTraceKeyValues(ctx)...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	sugar.Infow(fmt.Sprintf(format, args...), buildTraceKeyValues(ctx)...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	sugar.Warnw(fmt.Sprintf(format, args...), buildTraceKeyValues(ctx)...)
}

type options struct {
	m map[string]string
}
type Option func(o *options)

func MetaDataWithMap(m map[string]string) Option {
	return func(o *options) {
		o.m = m
	}
}
func Metadata(key, value string) Option {
	return func(o *options) {
		o.m[key] = value
	}
}

func AuditInfo(ctx context.Context, user, operator, operatorLevel string, request, response []byte, metadata ...Option) {
	sugar.Infow("", buildAuditKeyValues(ctx, user, operator, operatorLevel, request, response, metadata...)...)
}

func AuditError(ctx context.Context, user, operator, operatorLevel string, request, response []byte, metadata ...Option) {
	sugar.Errorw("", buildAuditKeyValues(ctx, user, operator, operatorLevel, request, response, metadata...)...)
}
func AuditDebug(ctx context.Context, user, operator, operatorLevel string, request, response []byte, metadata ...Option) {
	sugar.Debugw("", buildAuditKeyValues(ctx, user, operator, operatorLevel, request, response, metadata...)...)
}

func buildAuditKeyValues(ctx context.Context, user, operator, operatorLevel string, request, response []byte, metadata ...Option) []interface{} {
	res := buildTraceKeyValues(ctx)
	res = append(res, "user", user, "operator", operator, "request", request, "response", response, "operatorLevel", operatorLevel, "logtype", "audit")
	return res
}

func buildTraceKeyValues(ctx context.Context) []interface{} {
	res := make([]interface{}, 6)
	traceId := checkNil(ctx.Value(KeyTraceId))
	spanId := checkNil(ctx.Value(KeySpanId))
	parentId := checkNil(ctx.Value(KeyParentId))
	res[0] = KeyTraceId
	res[1] = traceId
	res[2] = KeyParentId
	res[3] = parentId
	res[4] = KeySpanId
	res[5] = spanId
	return res
}
func checkNil(arg interface{}) interface{} {
	if arg == nil {
		return ""
	}
	return arg
}
