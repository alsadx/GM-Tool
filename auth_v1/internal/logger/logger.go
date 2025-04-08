package logger

import (
	"context"
	"log/slog"
)

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	return ctx.Value("logger").(*slog.Logger)
}
