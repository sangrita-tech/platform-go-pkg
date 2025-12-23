package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg Config) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	encCfg := encoderConfigUTC()

	zcfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      cfg.DevMode,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    encCfg,
	}

	l, err := zcfg.Build()
	if err != nil {
		return nil, err
	}

	l = l.WithOptions(zap.AddCaller()).With(buildFields(cfg)...)
	return l, nil
}

func NewInMemory(cfg Config) (*zap.Logger, *Capture, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, nil, err
	}

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
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z07:00"))
		},
	}
}

func buildFields(cfg Config) []zap.Field {
	fields := []zap.Field{}
	for k, v := range cfg.BaseFields {
		fields = append(fields, zap.String(k, v))
	}
	return fields
}

func parseLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("logger: unknown level %q", level)
	}
}
