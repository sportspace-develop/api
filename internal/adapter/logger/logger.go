package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type option func(*zap.Logger)

func New(lvl string, options ...option) (*zap.Logger, error) {
	level := zap.AtomicLevel{}
	err := level.UnmarshalText([]byte(lvl))
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal level: %w", err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = level

	logger, err := cfg.Build()

	for _, opt := range options {
		opt(logger)
	}

	if err != nil {
		return nil, fmt.Errorf("failed create logger: %w", err)
	}
	return logger, err
}
