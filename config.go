package log

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

const (
	timestamp  = "timestamp"
	severity   = "severity"
	logger     = "logger"
	caller     = "caller"
	message    = "message"
	stacktrace = "stacktrace"
	jsonFormat = "json"
	textFormat = "text"
)

type Config struct {
	Application     string `mapstructure:"log_application"`
	Path            string `mapstructure:"log_path"`
	Level           string `mapstructure:"log_level"`
	Writers         string `mapstructure:"log_writers"`
	Format          string `mapstructure:"log_format"`
	Development     bool   `mapstructure:"log_development"`
	Color           bool   `mapstructure:"log_color"`
	RotateMaxBackup int    `mapstructure:"log_rotate_max_backup"`
	RotateMaxSize   int    `mapstructure:"log_rotate_max_size"`
	RotateMaxAge    int    `mapstructure:"log_rotate_max_age"`
	RotateCompress  bool   `mapstructure:"log_rotate_compress"`
}

func (cfg *Config) fileWriteSyncer(level string) zapcore.WriteSyncer {
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%v/%v.%v.log", cfg.Path, cfg.Application, level),
		MaxSize:    cfg.RotateMaxSize,
		MaxAge:     cfg.RotateMaxAge,
		MaxBackups: cfg.RotateMaxBackup,
		LocalTime:  true,
		Compress:   cfg.RotateCompress,
	})
	return fileWriteSyncer
}

func consoleEncoderConfig(enableColor bool) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       timestamp,
		LevelKey:      severity,
		NameKey:       logger,
		CallerKey:     caller,
		MessageKey:    message,
		StacktraceKey: stacktrace,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
			levelColor := levelSetColor(level)
			if enableColor {
				levelColor.EnableColor()
			}
			capitalString := levelColor.Sprintf(level.CapitalString())
			encoder.AppendString("[" + capitalString + "]")
		},
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: func(entryCaller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("[" + entryCaller.TrimmedPath() + "]")
		},
		ConsoleSeparator: " ",
	}
}

func jsonEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       timestamp,
		LevelKey:      severity,
		NameKey:       logger,
		CallerKey:     caller,
		MessageKey:    message,
		StacktraceKey: stacktrace,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(level.CapitalString())
		},
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: func(entryCaller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(entryCaller.TrimmedPath())
		},
	}
}

func mergeConfig(config, defaultConfig *Config) {
	if config.Application == "" {
		config.Application = defaultConfig.Application
	}
	if config.Path == "" {
		config.Path = defaultConfig.Path
	}
	if config.Level == "" {
		config.Level = defaultConfig.Level
	}
	if config.Writers == "" {
		config.Writers = defaultConfig.Writers
	}
	if config.Format == "" {
		config.Format = defaultConfig.Format
	}
	if config.RotateMaxBackup == 0 {
		config.RotateMaxBackup = defaultConfig.RotateMaxBackup
	}
	if config.RotateMaxAge == 0 {
		config.RotateMaxAge = defaultConfig.RotateMaxAge
	}
	if config.RotateMaxSize == 0 {
		config.RotateMaxSize = defaultConfig.RotateMaxSize
	}
}

func parseConfigFile(configFile string) (*Config, error) {
	config := &Config{}
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

var defaultConfig = &Config{
	Application:     "myapp",
	Path:            "logs",
	Level:           "info",
	Writers:         "stdout",
	Format:          "text",
	Development:     false,
	Color:           false,
	RotateMaxBackup: 10,
	RotateMaxSize:   10,
	RotateMaxAge:    7,
	RotateCompress:  false,
}
