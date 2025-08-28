package logging

import (
	"context"
	"errors"
	"flag"
	"io"
	stdlog "log"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
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
	zerologr.SetMaxV(level)

	var logger zerolog.Logger

	switch logFormat {
	case TextFormat:
		output := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: disableColor}
		output.TimeFormat = resolveTimestampFormat(timestampFormat)
		logger = zerolog.New(output).With().Timestamp().Caller().Logger()
	case JSONFormat:
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	default:
		return errors.New("log format unrecognized, pass `text` for text mode or `json` for JSON mode, passed : " + logFormat)
	}

	globalLog = zerologr.New(&logger)
	klog.SetLogger(globalLog.WithName("klog"))
	log.SetLogger(globalLog)

	// Set gRPC global logger
	grpcLogger := &grpcLogger{logger: logger}
	grpclog.SetLoggerV2(grpcLogger)

	return nil
}

func resolveTimestampFormat(format string) string {
	switch format {
	case ISO8601:
		return time.RFC3339
	case RFC3339:
		return time.RFC3339
	case MILLIS:
		return time.StampMilli
	case NANOS:
		return time.StampNano
	case EPOCH:
		return time.UnixDate
	case RFC3339NANO:
		return time.RFC3339Nano
	case DefaultTime:
		return time.RFC3339
	default:
		return time.RFC3339
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

func FromContext(ctx context.Context, keysAndValues ...interface{}) (logr.Logger, error) {
	logger, err := logr.FromContext(ctx)
	if err != nil {
		return logger, err
	}
	return logger.WithValues(keysAndValues...), nil
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
