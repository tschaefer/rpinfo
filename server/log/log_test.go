/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package log

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoggerReturnsErrorIfLevelIsUnknown(t *testing.T) {
	err := Logger("fatal", "structured")
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("unknown log level: fatal"), err)
}

func Test_LoggerReturnsErrorIfFormatIsUnknown(t *testing.T) {
	err := Logger("info", "xml")
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("unknown log format: xml"), err)
}

func Test_LoggerReturnsNoErrorIfLevelAndFormatAreKnown(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	formats := []string{"structured", "json"}

	for _, level := range levels {
		for _, format := range formats {
			err := Logger(level, format)
			assert.Nil(t, err)
		}
	}
}

func Test_RequestWritesLogMessage(t *testing.T) {
	var b strings.Builder
	w := io.Writer(&b)
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo})
	l := slog.New(h)
	slog.SetDefault(l)

	r := httptest.NewRequest("GET", "/temperature", nil)
	r.Header.Set("User-Agent", "rpinfo/1.0")
	Request(r, http.StatusOK, slog.LevelInfo, "This is a message")

	assert.Contains(t, b.String(), `"level":"INFO"`)
	assert.Contains(t, b.String(), `"msg":"This is a message"`)
	assert.Contains(t, b.String(), `"RemoteAddr":"`+r.RemoteAddr+`"`)
	assert.Contains(t, b.String(), `"UserAgent":"rpinfo/1.0"`)
	assert.Contains(t, b.String(), `"Status":200`)
	assert.Contains(t, b.String(), `"RequestMethod":"GET"`)
	assert.Contains(t, b.String(), `"RequestPath":"/temperature"`)
}

func Test_RequestDebugWritesDebugMessage(t *testing.T) {
	var b strings.Builder
	w := io.Writer(&b)
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := slog.New(h)
	slog.SetDefault(l)

	r := httptest.NewRequest("GET", "/temperature", nil)
	RequestDebug(r, http.StatusOK, "This is a debug message")

	assert.Contains(t, b.String(), `"level":"DEBUG"`)
	assert.Contains(t, b.String(), `"msg":"This is a debug message"`)
}

func Test_RequestErrorWritesErrorMessage(t *testing.T) {
	var b strings.Builder
	w := io.Writer(&b)
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelError})
	l := slog.New(h)
	slog.SetDefault(l)

	r := httptest.NewRequest("GET", "/temperature", nil)
	RequestError(r, http.StatusInternalServerError, "This is a error message")

	assert.Contains(t, b.String(), `"level":"ERROR"`)
	assert.Contains(t, b.String(), `"msg":"This is a error message"`)
}

func Test_RequestWarnWritesWarnMessage(t *testing.T) {
	var b strings.Builder
	w := io.Writer(&b)
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelWarn})
	l := slog.New(h)
	slog.SetDefault(l)

	r := httptest.NewRequest("GET", "/temperature", nil)
	RequestWarn(r, http.StatusInternalServerError, "This is a warning message")

	assert.Contains(t, b.String(), `"level":"WARN"`)
	assert.Contains(t, b.String(), `"msg":"This is a warning message"`)
}

func Test_RequestInfoWritesInfoMessage(t *testing.T) {
	var b strings.Builder
	w := io.Writer(&b)
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo})
	l := slog.New(h)
	slog.SetDefault(l)

	r := httptest.NewRequest("GET", "/temperature", nil)
	RequestInfo(r, http.StatusOK, "This is a info message")

	assert.Contains(t, b.String(), `"level":"INFO"`)
	assert.Contains(t, b.String(), `"msg":"This is a info message"`)
}
