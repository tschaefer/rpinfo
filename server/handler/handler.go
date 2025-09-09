/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package handler

import (
	"encoding/json"
	"iter"
	"log/slog"
	"maps"
	"net/http"
	"strings"

	"github.com/tschaefer/rpinfo/vcgencmd"
)

type Handle struct {
	Cmd vcgencmd.Exec
}

func runCmd(h Handle, w http.ResponseWriter, r *http.Request, args ...string) map[string]string {
	out, err := h.Cmd.Run(args...)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		go makeLog(r, http.StatusInternalServerError, slog.LevelError, err.Error())
		json.NewEncoder(w).Encode(map[string]string{"detail": "internal server error"})
		return nil
	}

	return out
}

func (h Handle) Temperature(w http.ResponseWriter, r *http.Request) {
	temp := runCmd(h, w, r, "measure_temp")
	if temp == nil {
		return
	}

	go makeLog(r, http.StatusOK, slog.LevelInfo, "Fetched temperature")
	json.NewEncoder(w).Encode(temp)
}

func (h Handle) Configuration(w http.ResponseWriter, r *http.Request) {
	options := []string{"int", "str"}
	config := make(map[string]string)
	for _, opt := range options {
		out := runCmd(h, w, r, "get_config", opt)
		if out == nil {
			return
		}

		maps.Copy(config, out)
	}

	go makeLog(r, http.StatusOK, slog.LevelInfo, "Fetched configuration")
	json.NewEncoder(w).Encode(config)
}

func (h Handle) Voltages(w http.ResponseWriter, r *http.Request) {
	options := []string{"core", "sdram_c", "sdram_i", "sdram_p"}
	voltages := make(map[string]string)
	for _, opt := range options {
		out := runCmd(h, w, r, "measure_volts", opt)
		if out == nil {
			return
		}

		voltages[opt] = out["volt"]
	}

	go makeLog(r, http.StatusOK, slog.LevelInfo, "Fetched voltages")
	json.NewEncoder(w).Encode(voltages)
}

func (h Handle) Throttled(w http.ResponseWriter, r *http.Request) {
	throttled := runCmd(h, w, r, "get_throttled")
	if throttled == nil {
		return
	}

	if r.URL.Query().Get("human") == "true" {
		messages, _ := parseThrottledHex(throttled["throttled"])
		message := strings.Join(messages, ", ")
		if len(message) == 0 {
			message = "No throttling"
		}
		throttled["throttled"] = message
	}

	go makeLog(r, http.StatusOK, slog.LevelInfo, "Fetched throttled status")
	json.NewEncoder(w).Encode(throttled)
}

func (h Handle) Clock(w http.ResponseWriter, r *http.Request) {
	options := []string{
		"arm", "core", "h264", "isp",
		"v3d", "uart", "pwm", "emmc",
		"pixel", "vec", "hdmi", "dpi",
	}
	clock := make(map[string]string)
	for _, opt := range options {
		out := runCmd(h, w, r, "measure_clock", opt)
		if out == nil {
			return
		}

		value := func() string {
			next, stop := iter.Pull(maps.Values(out))
			defer stop()

			v, _ := next()
			return v
		}

		clock[opt] = value()
	}

	go makeLog(r, http.StatusOK, slog.LevelInfo, "Fetched clock rates")
	json.NewEncoder(w).Encode(clock)
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
