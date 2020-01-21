package setting

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger custom logger for app
func InitLogger() *zap.Logger {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		TimeKey:        "time",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		CallerKey:      "caller",
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stderr),
		zapcore.DebugLevel,
	)
	return zap.New(core, zap.AddCaller())
}
