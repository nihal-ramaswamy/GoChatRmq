package utils

import (
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewZapLogger() *zap.Logger {
	env := GetDotEnvVariable(constants.ENV)

	switch env {
	case "release":
		return zap.Must(zap.NewProduction())
	case "debug":
		return zap.Must(zap.NewDevelopment())
	default:
		return zap.NewNop()
	}
}

func FxLogger() fx.Option {
	env := GetDotEnvVariable(constants.ENV)
	switch env {
	case "release":
		return fx.NopLogger
	case "debug":
		return fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		})
	default:
		return fx.NopLogger
	}
}
