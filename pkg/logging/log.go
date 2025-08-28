package logging

import (
	"context"
	"errors"
	"flag"
	"io"
	stdlog "log"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/grpclog"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	JSONFormat = "json"
	TextFormat = "text"
	LogLevel   = 1

	DefaultTime = "default"
	ISO8601     = "iso8601"
	RFC3339     = "rfc3339"
	MILLIS      = "millis"
	NANOS       = "nanos"
	EPOCH       = "epoch"
	RFC3339NANO = "rfc3339nano"
)

var globalLog = log.Log // Null log sink if SetLogger is not called

func InitFlags(flags *flag.FlagSet) {
	if flag.CommandLine.Lookup("log_dir") != nil {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}
	klog.InitFlags(flag.CommandLine)
}

func Setup(logFormat string, timestampFormat string, level int, disableColor bool) error {
	var logger *zap.Logger
	var err error

	switch logFormat {
	case TextFormat:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = resolveZapTimeEncoder(timestampFormat)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		if disableColor {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
		config.Level = zap.NewAtomicLevelAt(zapcore.Level(-level))
		logger, err = config.Build()
	case JSONFormat:
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = resolveZapTimeEncoder(timestampFormat)
		config.Level = zap.NewAtomicLevelAt(zapcore.Level(-level))
		logger, err = config.Build()
	default:
		return errors.New("log format unrecognized, pass `text` for text mode or `json` for JSON mode, passed : " + logFormat)
	}

	if err != nil {
		return err
	}

	// Create a zapr logger that implements logr.Logger
	zapLogger := zapr.NewLogger(logger)
	globalLog = zapLogger
	klog.SetLogger(globalLog.WithName("klog"))
	log.SetLogger(globalLog)

	// Set gRPC global logger
	grpcLogger := &grpcLogger{logger: logger}
	grpclog.SetLoggerV2(grpcLogger)

	return nil
}

func resolveZapTimeEncoder(format string) zapcore.TimeEncoder {
	switch format {
	case ISO8601, RFC3339, DefaultTime:
		return zapcore.ISO8601TimeEncoder
	case MILLIS:
		return zapcore.EpochMillisTimeEncoder
	case NANOS:
		return zapcore.EpochNanosTimeEncoder
	case EPOCH:
		return zapcore.EpochTimeEncoder
	case RFC3339NANO:
		return zapcore.RFC3339NanoTimeEncoder
	default:
		return zapcore.ISO8601TimeEncoder
	}
}

func GlobalLogger() logr.Logger {
	return globalLog
}

func WithName(name string) logr.Logger {
	return GlobalLogger().WithName(name)
}

func WithValues(keysAndValues ...interface{}) logr.Logger {
	return GlobalLogger().WithValues(keysAndValues...)
}

func V(level int) logr.Logger {
	return GlobalLogger().V(level)
}

func Info(msg string, keysAndValues ...interface{}) {
	GlobalLogger().WithCallDepth(1).Info(msg, keysAndValues...)
}

func Error(err error, msg string, keysAndValues ...interface{}) {
	GlobalLogger().WithCallDepth(1).Error(err, msg, keysAndValues...)
}

func FromContext(ctx context.Context) (logr.Logger, error) {
	return logr.FromContext(ctx)
}

func IntoContext(ctx context.Context, logger logr.Logger) context.Context {
	return logr.NewContext(ctx, logger)
}

func IntoBackground(log logr.Logger) context.Context {
	return IntoContext(context.Background(), log)
}

func IntoTODO(log logr.Logger) context.Context {
	return IntoContext(context.TODO(), log)
}

func Background() context.Context {
	return IntoBackground(GlobalLogger())
}

func TODO() context.Context {
	return IntoTODO(GlobalLogger())
}

type writerAdapter struct {
	io.Writer
	logger logr.Logger
}

func (w *writerAdapter) Write(p []byte) (int, error) {
	w.logger.Info(strings.TrimSuffix(string(p), "\n"))
	return len(p), nil
}

func StdLogger(log logr.Logger, prefix string) *stdlog.Logger {
	return stdlog.New(&writerAdapter{logger: log}, prefix, stdlog.LstdFlags)
}
