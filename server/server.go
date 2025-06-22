/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tschaefer/rpinfo/server/assets"
	"github.com/tschaefer/rpinfo/server/handler"
	"github.com/tschaefer/rpinfo/server/middleware"
	"github.com/tschaefer/rpinfo/vcgencmd"
	"github.com/tschaefer/rpinfo/version"
)

type Config struct {
	Port    string
	Host    string
	Auth    bool
	Token   string
	Metrics bool
	Redoc   bool
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
