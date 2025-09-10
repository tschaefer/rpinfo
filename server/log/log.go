/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package log

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func Logger(level, format string) error {
	var leveler slog.Leveler
	switch level {
	case "debug":
		leveler = slog.LevelDebug
	case "info":
		leveler = slog.LevelInfo
	case "warn":
		leveler = slog.LevelWarn
	case "error":
		leveler = slog.LevelError
	default:
		return fmt.Errorf("unknown log level: %s", level)
	}

	opts := &slog.HandlerOptions{
		Level: leveler,
	}

	var logger *slog.Logger
	switch format {
	case "structured":
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	case "text":
		// Use default logger, print info level only.
		return nil
	default:
		return fmt.Errorf("unknown log format: %s", format)
	}
	slog.SetDefault(logger)

	return nil
}

func Request(r *http.Request, status int, level slog.Level, msg string) {
	forwardedHeaders := []string{
		"X-Forwarded-For",
		"X-Real-IP",
	}
	remoteAddr := r.RemoteAddr
	for _, header := range forwardedHeaders {
		if ip := r.Header.Get(header); ip != "" {
			remoteAddr = ip
			break
		}
	}

	args := []any{
		slog.String("RemoteAddr", remoteAddr),
		slog.String("UserAgent", r.UserAgent()),
		slog.Int("Status", status),
		slog.String("RequestMethod", r.Method),
		slog.String("RequestPath", r.RequestURI),
	}

	switch level {
	case slog.LevelInfo:
		slog.Info(msg, args...)
	case slog.LevelWarn:
		slog.Warn(msg, args...)
	case slog.LevelError:
		slog.Error(msg, args...)
	default:
		slog.Info(msg, args...)
	}
}

func RequestInfo(r *http.Request, status int, msg string) {
	Request(r, status, slog.LevelInfo, msg)
}

func RequestWarn(r *http.Request, status int, msg string) {
	Request(r, status, slog.LevelWarn, msg)
}

func RequestError(r *http.Request, status int, msg string) {
	Request(r, status, slog.LevelError, msg)
}
