/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/

package handler

import (
	"bytes"
	"fmt"
	"iter"
	"maps"
	"net/http"
	"strconv"

	"github.com/VictoriaMetrics/metrics"
	"github.com/tschaefer/rpinfo/vcgencmd"
)

var (
	rpi = metrics.NewSet()
)

func Metrics(w http.ResponseWriter, r *http.Request) {
	clocks := []string{
		"arm", "core", "dpi", "emmc", "h264",
		"hdmi", "isp", "pixel", "pwm", "uart", "v3d", "vec",
	}
	for _, c := range clocks {
		name := fmt.Sprintf(`rpi_clock_%s`, c)
		rpi.GetOrCreateGauge(name, func() float64 { return clock(c) })
	}

	rpi.GetOrCreateGauge(`rpi_temperature`, func() float64 { return temperature() })

	voltages := []string{
		"core", "sdram_c", "sdram_i", "sdram_p",
	}
	for _, v := range voltages {
		name := fmt.Sprintf(`rpi_voltage_%s`, v)
		rpi.GetOrCreateGauge(name, func() float64 { return voltage(v) })
	}

	var buffer bytes.Buffer
	rpi.WritePrometheus(&buffer)
	if _, err := w.Write(buffer.Bytes()); err != nil {
		http.Error(w, "Failed to write metrics", http.StatusInternalServerError)
		return
	}
}

func clock(kind string) float64 {
	raw := exec("measure_clock", kind)
	if raw == nil {
		return 0.0
	}

	value := func() string {
		next, stop := iter.Pull(maps.Values(raw))
		defer stop()

		v, _ := next()
		return v
	}
	frequency, err := strconv.ParseFloat(value(), 64)
	if err != nil {
		return 0.0
	}

	return frequency
}

func temperature() float64 {
	raw := exec("measure_temp")
	if raw == nil {
		return 0.0
	}
	value := raw["temp"]
	temp, err := strconv.ParseFloat(value[:len(value)-2], 64)
	if err != nil {
		return 0.0
	}
	return temp
}

func voltage(kind string) float64 {
	raw := exec("measure_volts", kind)
	if raw == nil {
		return 0.0
	}
	value := raw["volt"]
	volt, err := strconv.ParseFloat(value[:len(value)-1], 64)
	if err != nil {
		return 0.0
	}
	return volt
}

func exec(args ...string) map[string]string {
	h := Handle{Cmd: vcgencmd.Cmd{}}
	return h.Cmd.Run(args...)
}
