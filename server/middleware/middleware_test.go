/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package middleware

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/tschaefer/rpinfo/version"
)

func Test_ResponseHeadersAreSet(t *testing.T) {
	version.GitCommit = "f0da3c4"

	req := httptest.NewRequest("GET", "/temperature", nil)
	rr := httptest.NewRecorder()

	handler := Headers(func(w http.ResponseWriter, r *http.Request) {})
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("X-Rpinfo-Commit") == "" {
		t.Errorf("Expected X-Rpinfo-Commit header to be set")
	}
	if rr.Header().Get("X-Rpinfo-Version") == "" {
		t.Errorf("Expected X-Rpinfo-Version header to be set")
	}
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be 'application/json', got %s", rr.Header().Get("Content-Type"))
	}
}

func Test_AuthorizationIsSkippedIfDisabled(t *testing.T) {
	req := httptest.NewRequest("GET", "/temperature", nil)
	rr := httptest.NewRecorder()

	handler := Authorization(false, "", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}
}

func Test_AuthorizationIsDeniedIfTokenIsInvalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/temperature", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	rr := httptest.NewRecorder()

	handler := Authorization(true, "valid_token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code 401, got %d", rr.Code)
	}
	if rr.Body.String() != "401 unauthorized\n" {
		t.Errorf("Expected body '401 unauthorized', got %s", rr.Body.String())
	}
}

func Test_AuthorizationIsDeniedIfNoHeaderPresent(t *testing.T) {
	req := httptest.NewRequest("GET", "/temperature", nil)
	rr := httptest.NewRecorder()

	handler := Authorization(true, "valid_token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code 401, got %d", rr.Code)
	}
	if rr.Body.String() != "401 unauthorized\n" {
		t.Errorf("Expected body '401 unauthorized', got %s", rr.Body.String())
	}
}

func Test_AuthorizationIsAcceptedIfTokenIsValid(t *testing.T) {
	req := httptest.NewRequest("GET", "/temperature", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	rr := httptest.NewRecorder()

	handler := Authorization(true, "valid_token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}
}
