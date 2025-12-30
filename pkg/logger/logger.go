package logger

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	klog "k8s.io/klog/v2"
)

func New(cfg *Config) (*zap.Logger, func(), error) {
	if cfg == nil {
		return nil, nil, errors.New("logger -> config is nil")
	}

	if err := cfg.validate(); err != nil {
		return nil, nil, fmt.Errorf("logger -> failed to validate config -> %w", err)
	}

	level := parseLevel(cfg.Level)
	encCfg := encoderConfigUTC()

	zcfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      cfg.DevMode,
		Encoding:         cfg.Format,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    encCfg,
	}

	l, err := zcfg.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("logger -> failed to build zap logger -> %w", err)
	}

	l = l.WithOptions(zap.AddCaller()).With(buildFields(cfg)...)

	restoreStdLog := zap.RedirectStdLog(l)

	klog.SetLogger(zapr.NewLogger(l))
	klog.EnableContextualLogging(true)

	cleanup := func() {
		restoreStdLog()
		_ = l.Sync()
	}

	return l, cleanup, nil
}

func NewInMemory(cfg *Config) (*zap.Logger, *Capture, error) {
	level := parseLevel(cfg.Level)
	encCfg := encoderConfigUTC()

	cap := &Capture{}
	encoder := zapcore.NewJSONEncoder(encCfg)
	core := zapcore.NewCore(encoder, zapcore.AddSync(cap), level)

	l := zap.New(core).WithOptions(zap.AddCaller()).With(buildFields(cfg)...)
	return l, cap, nil
}

func encoderConfigUTC() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		MessageKey:    "msg",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stack",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z07:00"))
		},
	}
}

func buildFields(cfg *Config) []zap.Field {
	fields := []zap.Field{}
	for k, v := range cfg.BaseFields {
		fields = append(fields, zap.String(k, v))
	}
	return fields
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
