package preset

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bamcop/kit/debug"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetDefaultSlog(filename string) {
	core := zapConsoleCore()
	if filename != "" {
		core = zapcore.NewTee(
			core,
			zapFileCore(filename),
		)
	}

	h := zapslog.HandlerOptions{
		AddSource: true,
	}.New(core)

	slog.SetDefault(slog.New(h))
}

func zapFileCore(filename string) zapcore.Core {
	if !(strings.HasPrefix(filename, "/") || strings.HasPrefix(filename, ".")) {
		filename = filepath.Join(debug.MustMainFileDir(), filename)
	}

	encoderConfig := zap.NewProductionEncoderConfig()

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    5,
		MaxBackups: 3,
		MaxAge:     3,
		Compress:   true,
	}

	writeSyncer := zapcore.AddSync(lumberJackLogger)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(writeSyncer),
		zap.InfoLevel,
	)

	return core
}

func zapConsoleCore() zapcore.Core {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		zap.InfoLevel,
	)

	return core
}
