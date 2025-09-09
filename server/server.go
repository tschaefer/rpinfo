/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package server

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/tschaefer/rpinfo/server/assets"
	"github.com/tschaefer/rpinfo/server/handler"
	"github.com/tschaefer/rpinfo/server/middleware"
	"github.com/tschaefer/rpinfo/vcgencmd"
	"github.com/tschaefer/rpinfo/version"
)

type Config struct {
	Port      string
	Host      string
	Auth      bool
	Token     string
	Metrics   bool
	Redoc     bool
	LogFormat string
	LogLevel  string
}

func Run(config Config) {
	Handler := handler.Handle{Cmd: vcgencmd.Cmd{}}

	router := mux.NewRouter()
	router.Handle("/temperature", middleware.ApplyAll(config.Auth, config.Token, Handler.Temperature)).Methods(http.MethodGet)
	router.Handle("/configuration", middleware.ApplyAll(config.Auth, config.Token, Handler.Configuration)).Methods(http.MethodGet)
	router.Handle("/voltages", middleware.ApplyAll(config.Auth, config.Token, Handler.Voltages)).Methods(http.MethodGet)
	router.Handle("/throttled", middleware.ApplyAll(config.Auth, config.Token, Handler.Throttled)).Methods(http.MethodGet)
	router.Handle("/clock", middleware.ApplyAll(config.Auth, config.Token, Handler.Clock)).Methods(http.MethodGet)

	if config.Redoc {
		router.PathPrefix("/redoc").Handler(http.StripPrefix("/redoc", http.FileServer(http.FS(assets.StaticContent))))
	}

	if config.Metrics {
		router.HandleFunc("/metrics", middleware.Authorization(config.Auth, config.Token, handler.Metrics)).Methods(http.MethodGet)
	}

	router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handler.MethodNotAllowedHandler)

	if err := setLogger(config); err != nil {
		log.Fatalf("Failed to set logger: %v", err)
	}

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", config.Host, config.Port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		Handler:        router,
	}

	log.Printf("Starting rpinfo server. Version: %s - %s", version.Release(), version.Commit())
	log.Printf("Listening on %s:%s, auth: %t, metrics: %t, redoc: %t", config.Host, config.Port, config.Auth, config.Metrics, config.Redoc)
	log.Fatal(server.ListenAndServe())
}

func setLogger(config Config) error {
	var leveler slog.Leveler
	switch config.LogLevel {
	case "debug":
		leveler = slog.LevelDebug
	case "info":
		leveler = slog.LevelInfo
	case "warn":
		leveler = slog.LevelWarn
	case "error":
		leveler = slog.LevelError
	default:
		return fmt.Errorf("unknown log level: %s", config.LogLevel)
	}

	opts := &slog.HandlerOptions{
		Level: leveler,
	}

	var logger *slog.Logger
	switch config.LogFormat {
	case "structured":
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	case "text":
		// Use default logger, print info level only.
		return nil
	default:
		return fmt.Errorf("unknown log format: %s", config.LogFormat)
	}
	slog.SetDefault(logger)

	return nil
}
