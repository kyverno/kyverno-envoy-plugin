package logging

import "go.uber.org/zap"

type grpcLogger struct {
	logger *zap.Logger
}

func (g *grpcLogger) Info(args ...interface{}) {
	g.logger.Info("", zap.Any("args", args))
}

func (g *grpcLogger) Infoln(args ...interface{}) {
	g.logger.Info("", zap.Any("args", args))
}

func (g *grpcLogger) Infof(format string, args ...interface{}) {
	g.logger.Sugar().Infof(format, args...)
}

func (g *grpcLogger) Warning(args ...interface{}) {
	g.logger.Warn("", zap.Any("args", args))
}

func (g *grpcLogger) Warningln(args ...interface{}) {
	g.logger.Warn("", zap.Any("args", args))
}

func (g *grpcLogger) Warningf(format string, args ...interface{}) {
	g.logger.Sugar().Warnf(format, args...)
}

func (g *grpcLogger) Error(args ...interface{}) {
	g.logger.Error("", zap.Any("args", args))
}

func (g *grpcLogger) Errorln(args ...interface{}) {
	g.logger.Error("", zap.Any("args", args))
}

func (g *grpcLogger) Errorf(format string, args ...interface{}) {
	g.logger.Sugar().Errorf(format, args...)
}

func (g *grpcLogger) Fatal(args ...interface{}) {
	g.logger.Fatal("", zap.Any("args", args))
}

func (g *grpcLogger) Fatalln(args ...interface{}) {
	g.logger.Fatal("", zap.Any("args", args))
}

func (g *grpcLogger) Fatalf(format string, args ...interface{}) {
	g.logger.Sugar().Fatalf(format, args...)
}

func (g *grpcLogger) V(level int) bool {
	// gRPC uses levels 0-99, where 0 is most verbose
	// We'll map this to zap levels (0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug, 6=debug+1)
	// For gRPC level 0-2, we'll return true (allow logging), for higher levels we'll check against our configured level
	return level <= 2 || level <= LogLevel
}
