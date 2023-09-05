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
	"net/http"
	"net/http/httptest"
	"testing"
)

//go:embed html/headers.html
var headersHTML string

func headersBody(headers http.Header) string {
	var body bytes.Buffer
	tmpl := template.Must(template.New("test").Parse(headersHTML))
	tmpl.Execute(&body, HeadersPageData{
		Title:   "Request Headers",
		Headers: NewHeaderInfo(headers),
	})

	return body.String()
}

// TestheadersHandler tests the Headers handler.
func TestHeadersHandler(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
	}{
		{
			name: "Empty headers",
		},
		{
			name: "Multiple headers",
			headers: http.Header{
				"Content-Type":    {"application/json"},
				"X-Custom-Header": {"value"},
				"Accept-Encoding": {"gzip"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			req.Header = tt.headers
			handler.HeadersHandler(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			body := rr.Body.String()
			expectedBody := headersBody(tt.headers)

			if body != expectedBody {
				t.Errorf("expected response body '%s', got '%s'", expectedBody, rr.Body.String())
			}
		})
	}
}
