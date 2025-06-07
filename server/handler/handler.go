/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package handler

import (
	"encoding/json"
	"iter"
	"maps"
	"net/http"
	"strings"

	"github.com/tschaefer/rpinfo/vcgencmd"
)

type Handle struct {
	Cmd vcgencmd.Exec
}

func runCmd(h Handle, w http.ResponseWriter, args ...string) map[string]string {
	out := h.Cmd.Run(args...)
	if out == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"detail": "internal server error"})
		return nil
	}

	return out
}

func (h Handle) Temperature(w http.ResponseWriter, r *http.Request) {
	temp := runCmd(h, w, "measure_temp")
	if temp == nil {
		return
	}

	json.NewEncoder(w).Encode(temp)
}

func (h Handle) Configuration(w http.ResponseWriter, r *http.Request) {
	options := []string{"int", "str"}
	config := make(map[string]string)
	for _, opt := range options {
		out := runCmd(h, w, "get_config", opt)
		if out == nil {
			return
		}

		maps.Copy(config, out)
	}

	json.NewEncoder(w).Encode(config)
}

func (h Handle) Voltages(w http.ResponseWriter, r *http.Request) {
	options := []string{"core", "sdram_c", "sdram_i", "sdram_p"}
	voltages := make(map[string]string)
	for _, opt := range options {
		out := runCmd(h, w, "measure_volts", opt)
		if out == nil {
			return
		}

		voltages[opt] = out["volt"]
	}

	json.NewEncoder(w).Encode(voltages)
}

func (h Handle) Throttled(w http.ResponseWriter, r *http.Request) {
	throttled := runCmd(h, w, "get_throttled")
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
		out := runCmd(h, w, "measure_clock", opt)
		if out == nil {
			return
		}

		next, stop := iter.Pull(maps.Values(out))
		defer stop()

		clock[opt], _ = next()
	}

	json.NewEncoder(w).Encode(clock)
}
