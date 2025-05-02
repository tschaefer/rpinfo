/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package middleware

import (
	"net/http"
	"strings"

	"github.com/tschaefer/rpinfo/version"
)

func Headers(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Rpinfo-Commit", version.Commit())
		w.Header().Set("X-Rpinfo-Version", version.Release())
		w.Header().Set("Content-Type", "application/json")

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
		parts := strings.SplitN(bearer, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] != token {
			http.Error(w, "401 unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func ApplyAll(auth bool, token string, next http.HandlerFunc) http.HandlerFunc {
	// middleware is applied in reverse order
	next = Headers(next)
	next = Authorization(auth, token, next)

	return next
}
