package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

func NewLogger(config *Config) *zap.SugaredLogger {
	return newLogger(config)
}

func NewLoggerWithConfigFile(configFile string) *zap.SugaredLogger {
	config, err := parseConfigFile(configFile)
	if err != nil {
		fmt.Println("parse config error, using default config")
		config = defaultConfig
	}
	return newLogger(config)
}

func newLogger(config *Config) *zap.SugaredLogger {
	mergeConfig(config, defaultConfig)

	debugPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.DebugLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel
	})
	warnPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.WarnLevel
	})
	errorPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})

	var cores []zapcore.Core
	var encoder zapcore.Encoder
	switch config.Format {
	case jsonFormat:
		encoder = zapcore.NewJSONEncoder(jsonEncoderConfig())
	case textFormat:
		encoder = zapcore.NewConsoleEncoder(consoleEncoderConfig(false))
	}

	for _, v := range strings.Split(config.Writers, ",") {
		switch {
		case v == "stdout":
			if config.Color {
				colorEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig(true))
				cores = append(
					cores,
					zapcore.NewCore(colorEncoder, zapcore.AddSync(os.Stdout), debugPriority),
					zapcore.NewCore(colorEncoder, zapcore.AddSync(os.Stdout), infoPriority),
					zapcore.NewCore(colorEncoder, zapcore.AddSync(os.Stdout), warnPriority),
					zapcore.NewCore(colorEncoder, zapcore.AddSync(os.Stdout), errorPriority),
				)
			} else {
				cores = append(
					cores,
					zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), debugPriority),
					zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), infoPriority),
					zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), warnPriority),
					zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), errorPriority),
				)
			}
		case v == "file":
			cores = append(
				cores,
				zapcore.NewCore(encoder, config.fileWriteSyncer("debug"), debugPriority),
				zapcore.NewCore(encoder, config.fileWriteSyncer("info"), infoPriority),
				zapcore.NewCore(encoder, config.fileWriteSyncer("warn"), warnPriority),
				zapcore.NewCore(encoder, config.fileWriteSyncer("error"), errorPriority),
			)
		}
	}

	logger := zap.New(zapcore.NewTee(cores...)).WithOptions(
		zap.AddCaller(),
		zap.AddStacktrace(errorPriority),
	)

	if config.Development {
		logger = logger.WithOptions(zap.Development())
	}

	return logger.Sugar()
}
