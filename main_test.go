/*
Copyright 2023 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// handler to use across all tests.
var handler *Handler

const appName = "Test App Name"

// TestMain sets up the logger and then invocates all the tests.
func TestMain(m *testing.M) {
	// configure logger
	opts := &slog.HandlerOptions{
		//Level:     slog.LevelDebug,
		AddSource: true,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, opts))
	slog.SetDefault(logger)

	// initialize templates
	tmpl, err := InitTemplates("html/*.html")
	if err != nil {
		slog.Error("failed to InitTemplates", "err", err)
	}

	handler = NewHandler(appName, tmpl)

	os.Exit(m.Run())
}

//go:embed html/root.html
var rootHTML string

// TestRootHandler tests the Root handler.
func TestRootHandler(t *testing.T) {
	var goodBody bytes.Buffer
	template.Must(template.New("test").Parse(rootHTML)).Execute(&goodBody, RootPageData{Title: appName})

	tests := []struct {
		name             string
		path             string
		method           string
		expectedCode     int
		expectedBody     string
		expectedLocation string
	}{
		{
			name:         "Root path",
			path:         "/",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: goodBody.String(),
		},
		{
			name:         "Non-root path",
			path:         "/other",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
			expectedBody: "404 page not found",
		},
		{
			name:         "Invalid method",
			path:         "/",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: http.MethodPut + " " + http.StatusText(http.StatusMethodNotAllowed),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler.RootHandler(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			if tt.expectedLocation != "" {
				location := rr.Header().Get("Location")
				if location != tt.expectedLocation {
					t.Errorf("expected Location header '%s', got '%s'", tt.expectedLocation, location)
				}
			}

			body := strings.TrimSpace(rr.Body.String())
			if body != tt.expectedBody {
				t.Errorf("expected response body '%s', got '%s'", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestHelloHandler(t *testing.T) {
	testCases := []struct {
		name             string
		method           string
		expectedStatus   int
		expectedBody     string
		expectedCacheCtl string
	}{
		{
			name:             "GET request",
			method:           http.MethodGet,
			expectedStatus:   http.StatusOK,
			expectedBody:     "hello",
			expectedCacheCtl: "no-cache, no-store, must-revalidate",
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   http.MethodPost + " " + http.StatusText(http.StatusMethodNotAllowed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "/hello", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.HelloHandler(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			body := strings.TrimSpace(rr.Body.String())
			if body != tc.expectedBody {
				t.Errorf("expected response body '%s', got '%s'", tc.expectedBody, body)
			}

			cacheCtl := rr.Header().Get("Cache-Control")
			if cacheCtl != tc.expectedCacheCtl {
				t.Errorf("expected Cache-Control header '%s', got '%s'", tc.expectedCacheCtl, cacheCtl)
			}
		})
	}
}

func TestIPHandler(t *testing.T) {
	testCases := []struct {
		name           string
		request        *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid request",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/ip", nil)
				req.RemoteAddr = "127.0.0.1:1234"
				return req
			}(),
			expectedStatus: http.StatusOK,
			expectedBody:   "RemoteAddr: 127.0.0.1:1234",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler.RemoteHandler(rr, tc.request)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			body := strings.TrimSpace(rr.Body.String())
			if body != tc.expectedBody {
				t.Errorf("expected response body '%s', got '%s'", tc.expectedBody, body)
			}
		})
	}
}

func TestRequestHandler(t *testing.T) {
	testCases := []struct {
		name           string
		request        *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid request",
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/request", nil)
				return req
			}(),
			expectedStatus: http.StatusOK,
			expectedBody:   "GET /request HTTP/1.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler.RequestHandler(rr, tc.request)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			body := strings.TrimSpace(rr.Body.String())
			if body != tc.expectedBody {
				t.Errorf("expected response body '%s', got '%s'", tc.expectedBody, body)
			}
		})
	}
}
