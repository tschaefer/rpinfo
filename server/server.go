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

func Run(port string, host string, auth bool, token string) {
	Handler := handler.Handle{Cmd: vcgencmd.Cmd{}}

	router := mux.NewRouter()
	router.Handle("/temperature", middleware.ApplyAll(auth, token, Handler.Temperature)).Methods(http.MethodGet)
	router.Handle("/configuration", middleware.ApplyAll(auth, token, Handler.Configuration)).Methods(http.MethodGet)
	router.Handle("/voltages", middleware.ApplyAll(auth, token, Handler.Voltages)).Methods(http.MethodGet)
	router.Handle("/throttled", middleware.ApplyAll(auth, token, Handler.Throttled)).Methods(http.MethodGet)
	router.PathPrefix("/redoc").Handler(http.StripPrefix("/redoc", http.FileServer(http.FS(assets.StaticContent))))

	router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handler.MethodNotAllowedHandler)

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", host, port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		Handler:        router,
	}

	log.Printf("Starting rpinfo server. Version: %s - %s", version.Release(), version.Commit())
	log.Printf("Listening on %s:%s, auth: %t", host, port, auth)
	log.Fatal(server.ListenAndServe())
}
