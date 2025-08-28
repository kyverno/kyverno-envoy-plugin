package logging

import "github.com/rs/zerolog"

type grpcLogger struct {
	logger zerolog.Logger
}

func (g *grpcLogger) Info(args ...interface{}) {
	g.logger.Info().Msgf("%v", args...)
}

func (g *grpcLogger) Infoln(args ...interface{}) {
	g.logger.Info().Msgf("%v", args...)
}

func (g *grpcLogger) Infof(format string, args ...interface{}) {
	g.logger.Info().Msgf(format, args...)
}

func (g *grpcLogger) Warning(args ...interface{}) {
	g.logger.Warn().Msgf("%v", args...)
}

func (g *grpcLogger) Warningln(args ...interface{}) {
	g.logger.Warn().Msgf("%v", args...)
}

func (g *grpcLogger) Warningf(format string, args ...interface{}) {
	g.logger.Warn().Msgf(format, args...)
}

func (g *grpcLogger) Error(args ...interface{}) {
	g.logger.Error().Msgf("%v", args...)
}

func (g *grpcLogger) Errorln(args ...interface{}) {
	g.logger.Error().Msgf("%v", args...)
}

func (g *grpcLogger) Errorf(format string, args ...interface{}) {
	g.logger.Error().Msgf(format, args...)
}

func (g *grpcLogger) Fatal(args ...interface{}) {
	g.logger.Fatal().Msgf("%v", args...)
}

func (g *grpcLogger) Fatalln(args ...interface{}) {
	g.logger.Fatal().Msgf("%v", args...)
}

func (g *grpcLogger) Fatalf(format string, args ...interface{}) {
	g.logger.Fatal().Msgf(format, args...)
}

func (g *grpcLogger) V(level int) bool {
	// gRPC uses levels 0-99, where 0 is most verbose
	// We'll map this to zerolog levels (0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug, 6=trace)
	// For gRPC level 0-2, we'll return true (allow logging), for higher levels we'll check against our configured level
	return level <= 2 || level <= LogLevel
}
