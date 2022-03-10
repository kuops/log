package log

import (
	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"
)

func levelSetColor(level zapcore.Level) *color.Color {
	levelColor := &color.Color{}
	switch {
	case level == zapcore.DebugLevel:
		levelColor = color.New(color.FgMagenta)
	case level == zapcore.InfoLevel:
		levelColor = color.New(color.FgCyan)
	case level == zapcore.WarnLevel:
		levelColor = color.New(color.FgYellow)
	case level >= zapcore.ErrorLevel:
		levelColor = color.New(color.FgRed)
	}
	return levelColor
}
