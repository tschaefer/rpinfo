/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package handler

import (
	"encoding/json"
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
		http.Error(w, "500 internal server error", http.StatusInternalServerError)

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

		for k, v := range out {
			config[k] = v
		}
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
