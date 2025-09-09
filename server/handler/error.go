/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/tschaefer/rpinfo/version"
)

func JSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Rpinfo-Commit", version.Commit())
	w.Header().Set("X-Rpinfo-Version", version.Release())
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"detail": message})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	go makeLog(r, http.StatusInternalServerError, slog.LevelError, "not found")
	JSONError(w, http.StatusNotFound, "not found")
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	go makeLog(r, http.StatusMethodNotAllowed, slog.LevelWarn, "method not allowed")
	JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
}
