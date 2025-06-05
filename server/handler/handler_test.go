/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tschaefer/rpinfo/server/assets"
)

type mockRunnerSuccess struct{}

func (m mockRunnerSuccess) Run(args ...string) map[string]string {
	switch args[0] {
	case "measure_temp":
		return map[string]string{"temp": "45.0'C"}
	case "measure_volts":
		switch args[1] {
		case "core":
			return map[string]string{"volt": "1.3500V"}
		case "sdram_c":
			return map[string]string{"volt": "1.2000V"}
		case "sdram_i":
			return map[string]string{"volt": "1.2000V"}
		case "sdram_p":
			return map[string]string{"volt": "1.2250V"}
		default:
			return nil
		}
	case "get_config":
		return map[string]string{"init_uart_clock": "0x2dc6c00", "overlay_prefix": "overlays/", "total_mem": "512"}
	case "get_throttled":
		return map[string]string{"throttled": "0x50000"}
	default:
		return nil
	}
}

type mockRunnerError struct{}

func (m mockRunnerError) Run(args ...string) map[string]string {
	return nil
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
