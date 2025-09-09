/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/tschaefer/rpinfo/version"
)

func JSONError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"detail": message})
}

func ResponseHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Rpinfo-Commit", version.Commit())
		w.Header().Set("X-Rpinfo-Version", version.Release())
		w.Header().Set("Content-Type", "application/json")

		next(w, r)
	}
}

func RequestHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		if accept == "" || (accept != "application/json" && accept != "*/*") {
			go makeLog(r, http.StatusNotAcceptable, slog.LevelWarn, "not acceptable")
			JSONError(w, http.StatusNotAcceptable, "not acceptable")
			return
		}
		next(w, r)
	}
}

func Authorization(auth bool, token string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth {
			next(w, r)
			return
		}

		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			go makeLog(r, http.StatusUnauthorized, slog.LevelWarn, "unauthorized")
			JSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		parts := strings.SplitN(bearer, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] != token {
			go makeLog(r, http.StatusForbidden, slog.LevelWarn, "forbidden")
			JSONError(w, http.StatusForbidden, "forbidden")
			return
		}

		next(w, r)
	}
}

func ApplyAll(auth bool, token string, next http.HandlerFunc) http.HandlerFunc {
	// middleware is applied in reverse order
	next = RequestHeaders(next)
	next = Authorization(auth, token, next)
	next = ResponseHeaders(next)

	return next
}

func makeLog(r *http.Request, status int, level slog.Level, msg string) {
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
