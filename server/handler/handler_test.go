/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tschaefer/rpinfo/server/assets"
)

type mockRunnerSuccess struct{}

func (m mockRunnerSuccess) Run(args ...string) (map[string]string, error) {
	switch args[0] {
	case "measure_temp":
		return map[string]string{"temp": "45.0'C"}, nil
	case "measure_volts":
		switch args[1] {
		case "core":
			return map[string]string{"volt": "1.3500V"}, nil
		case "sdram_c":
			return map[string]string{"volt": "1.2000V"}, nil
		case "sdram_i":
			return map[string]string{"volt": "1.2000V"}, nil
		case "sdram_p":
			return map[string]string{"volt": "1.2250V"}, nil
		default:
			return nil, nil
		}
	case "get_config":
		return map[string]string{"init_uart_clock": "0x2dc6c00", "overlay_prefix": "overlays/", "total_mem": "512"}, nil
	case "get_throttled":
		return map[string]string{"throttled": "0x50000"}, nil
	case "measure_clock":
		switch args[1] {
		case "arm":
			return map[string]string{"freq": "600000000"}, nil
		case "core":
			return map[string]string{"freq": "250000000"}, nil
		default:
			return map[string]string{"freq": "0"}, nil
		}
	default:
		return nil, nil
	}
}

type mockRunnerError struct{}

func (m mockRunnerError) Run(args ...string) (map[string]string, error) {
	return nil, fmt.Errorf("command failed")
}

func Test_TemperatureReturnsJSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/temperature", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerSuccess{}}
	handler := http.HandlerFunc(Handler.Temperature)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"temp":"45.0'C"}`
	got := rr.Body.String()
	got = strings.TrimSpace(got)
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_TemperatureReturnsServerErrorIfCommandFails(t *testing.T) {
	req := httptest.NewRequest("GET", "/temperature", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerError{}}
	handler := http.HandlerFunc(Handler.Temperature)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func Test_VoltagesReturnsJSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/voltages", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerSuccess{}}
	handler := http.HandlerFunc(Handler.Voltages)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"core":"1.3500V","sdram_c":"1.2000V","sdram_i":"1.2000V","sdram_p":"1.2250V"}`
	got := rr.Body.String()
	got = strings.TrimSpace(got)
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_VoltagesReturnsServerErrorIfCommandFails(t *testing.T) {
	req := httptest.NewRequest("GET", "/voltages", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerError{}}
	handler := http.HandlerFunc(Handler.Voltages)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func Test_ConfigurationReturnsJSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/configuration", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerSuccess{}}
	handler := http.HandlerFunc(Handler.Configuration)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"init_uart_clock":"0x2dc6c00","overlay_prefix":"overlays/","total_mem":"512"}`
	got := rr.Body.String()
	got = strings.TrimSpace(got)
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_ConfigurationReturnsServerErrorIfCommandFails(t *testing.T) {
	req := httptest.NewRequest("GET", "/configuration", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerError{}}
	handler := http.HandlerFunc(Handler.Configuration)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func Test_ThrottledReturnsJSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/throttled", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerSuccess{}}
	handler := http.HandlerFunc(Handler.Throttled)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"throttled":"0x50000"}`
	got := rr.Body.String()
	got = strings.TrimSpace(got)
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	req = httptest.NewRequest("GET", "/throttled?human=true", nil)
	rr = httptest.NewRecorder()

	Handler = Handle{Cmd: mockRunnerSuccess{}}
	handler = http.HandlerFunc(Handler.Throttled)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected = `{"throttled":"Undervoltage has occurred, Throttling has occurred"}`
	got = rr.Body.String()
	got = strings.TrimSpace(got)
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_ThrottledReturnsServerErrorIfCommandFails(t *testing.T) {
	req := httptest.NewRequest("GET", "/throttled", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerError{}}
	handler := http.HandlerFunc(Handler.Throttled)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func Test_RedocReturnsHTML(t *testing.T) {
	req := httptest.NewRequest("GET", "/redoc", nil)
	rr := httptest.NewRecorder()
	handler := http.Handler(http.StripPrefix("/redoc", http.FileServer(http.FS(assets.StaticContent))))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "<!DOCTYPE html>"
	if !strings.HasPrefix(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want prefix %v",
			rr.Body.String(), expected)
	}

	expected = "text/html; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expected)
	}
}

func Test_MeasureClockReturnsJSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/measure_clock", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerSuccess{}}
	handler := http.HandlerFunc(Handler.Clock)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"arm":"600000000","core":"250000000","dpi":"0","emmc":"0","h264":"0","hdmi":"0","isp":"0","pixel":"0","pwm":"0","uart":"0","v3d":"0","vec":"0"}`
	got := rr.Body.String()
	got = strings.TrimSpace(got)
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_MeasureClockReturnsServerErrorIfCommandFails(t *testing.T) {
	req := httptest.NewRequest("GET", "/measure_clock", nil)
	rr := httptest.NewRecorder()

	Handler := Handle{Cmd: mockRunnerError{}}
	handler := http.HandlerFunc(Handler.Clock)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func Test_MetricsReturnsPrometheusText(t *testing.T) {
	req := httptest.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Metrics)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "rpi_clock_arm 0\nrpi_clock_core 0\nrpi_clock_dpi 0\n" +
		"rpi_clock_emmc 0\nrpi_clock_h264 0\nrpi_clock_hdmi 0\n" +
		"rpi_clock_isp 0\nrpi_clock_pixel 0\nrpi_clock_pwm 0\n" +
		"rpi_clock_uart 0\nrpi_clock_v3d 0\nrpi_clock_vec 0\n" +
		"rpi_temperature 0\nrpi_voltage_core 0\nrpi_voltage_sdram_c 0\n" +
		"rpi_voltage_sdram_i 0\nrpi_voltage_sdram_p 0"

	got := rr.Body.String()
	got = strings.TrimSpace(got)
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
