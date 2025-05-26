package logs

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities"
	sloghttp "github.com/samber/slog-http"
)

func Init(serviceName string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(handler).With("service", serviceName)
	return logger
}

func Fatal(logger *slog.Logger, msg string, attrs ...any) {
	logger.Error(msg, attrs...)
	os.Exit(1)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	options := sloghttp.Config{
		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,
	}
	return func(next http.Handler) http.Handler {
		h := sloghttp.Recovery(next)
		h = utilities.RequestIDMiddleware(h)
		h = sloghttp.NewWithConfig(logger, options)(h)
		return h
	}
}
