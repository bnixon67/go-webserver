package main

import (
	_ "embed"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"
)

//go:embed html/root.html
var helloHtml string

// TestMain sets up the logger and then invocates all the tests.
func TestMain(m *testing.M) {
	// configure logger
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, opts))
	slog.SetDefault(logger)

	os.Exit(m.Run())
}

// TestRootHandler tests the Root handler.
func TestRootHandler(t *testing.T) {
	handler := Handler{}

	tests := []struct {
		name             string
		path             string
		expectedCode     int
		expectedBody     string
		expectedLocation string
	}{
		{
			name:         "Root path",
			path:         "/",
			expectedCode: http.StatusOK,
			expectedBody: strings.TrimSpace(helloHtml),
		},
		{
			name:         "Non-root path",
			path:         "/other",
			expectedCode: http.StatusNotFound,
			expectedBody: "404 page not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.Root(rr, req)

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

// TestheadersHandler tests the Headers handler.
func TestHeadersHandler(t *testing.T) {
	handler := Handler{}

	tests := []struct {
		name         string
		headers      map[string]string
		expectedKeys []string
	}{
		{
			name: "Empty headers",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			expectedKeys: []string{"Content-Type"},
		},
		{
			name: "Multiple headers",
			headers: map[string]string{
				"Content-Type":    "application/json",
				"X-Custom-Header": "value",
				"Accept-Encoding": "gzip",
			},
			expectedKeys: []string{"Accept-Encoding", "Content-Type", "X-Custom-Header"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			for key, value := range tt.headers {
				req.Header.Add(key, value)
			}

			rr := httptest.NewRecorder()
			handler.Headers(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			response := rr.Body.String()
			keys := extractKeysFromResponse(response)

			if !stringSlicesEqual(tt.expectedKeys, keys) {
				t.Errorf("expected headers keys %v, got %v", tt.expectedKeys, keys)
			}
		})
	}
}

func extractKeysFromResponse(response string) []string {
	lines := strings.Split(response, "\n")
	keys := make([]string, 0)

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			keys = append(keys, strings.TrimSpace(parts[0]))
		}
	}

	sort.Strings(keys)
	return keys
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
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
		// Add more test cases as needed
	}

	handler := Handler{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "/hello", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.Hello(rr, req)

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
		// Add more test cases as needed
	}

	handler := Handler{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler.RemoteAddr(rr, tc.request)

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
		// Add more test cases as needed
	}

	handler := Handler{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler.Request(rr, tc.request)

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
