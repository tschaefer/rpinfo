/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tschaefer/rpinfo/server/log"
	"github.com/tschaefer/rpinfo/version"
)

func JSONError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"detail": message})
}

func ResponseHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Rpinfo-Commit", version.Commit())
		w.Header().Set("X-Rpinfo-Version", version.Release())
		w.Header().Set("Content-Type", "application/json")

		next(w, r)
	}
}

func RequestHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		if accept == "" || (accept != "application/json" && accept != "*/*") {
			go log.RequestWarn(r, http.StatusNotAcceptable, "not acceptable")
			JSONError(w, http.StatusNotAcceptable, "not acceptable")
			return
		}
		next(w, r)
	}
}

func Authorization(auth bool, token string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth {
			next(w, r)
			return
		}

		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			go log.RequestWarn(r, http.StatusUnauthorized, "unauthorized")
			JSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		parts := strings.SplitN(bearer, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] != token {
			go log.RequestWarn(r, http.StatusForbidden, "forbidden")
			JSONError(w, http.StatusForbidden, "forbidden")
			return
		}

		next(w, r)
	}
}

func ApplyAll(auth bool, token string, next http.HandlerFunc) http.HandlerFunc {
	// middleware is applied in reverse order
	next = RequestHeaders(next)
	next = Authorization(auth, token, next)
	next = ResponseHeaders(next)

	return next
}
