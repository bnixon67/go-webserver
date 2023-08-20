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
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"sort"
	"strings"
)

// Handler is responsible for handling HTTP requests.
type Handler struct{}

//go:embed html/root.html
var rootHTML string

// Root handles the root ("/") route.
func (h Handler) Root(w http.ResponseWriter, r *http.Request) {
	logger := Logger(r.Context())
	logger.Debug("Root Handler")

	if r.URL.Path != "/" {
		logger.Warn("invalid path")
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprint(w, rootHTML)
}

// Hello responds with a simple "hello" message.
func (h Handler) Hello(w http.ResponseWriter, r *http.Request) {
	logger := Logger(r.Context())
	logger.Debug("Hello Handler")

	// try and force client not to cache content
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Fprintln(w, "hello")
}

// Headers prints the headers of the request in sorted order.
func (h Handler) Headers(w http.ResponseWriter, r *http.Request) {
	logger := Logger(r.Context())
	logger.Debug("Headers Handler")

	// get header keys
	keys := make([]string, 0, len(r.Header))
	for key := range r.Header {
		keys = append(keys, key)
	}

	// sort headers
	sort.Strings(keys)

	// print key-value pairs
	for _, key := range keys {
		value := strings.Join(r.Header[key], ", ")
		fmt.Fprintf(w, "%v: %v\n", key, value)
	}
}

// RemoteAddr responds with the RemoteAddr and common headers for the actual RemoteAddr.
// Note: RemoteAddr may not be valid if a proxy, load balancer, or similar is used to route the request.
func (h Handler) RemoteAddr(w http.ResponseWriter, r *http.Request) {
	logger := Logger(r.Context())
	logger.Debug("Headers Handler")

	// show values
	fmt.Fprintf(w, "RemoteAddr: %v\n", r.RemoteAddr)

	headers := []string{
		"Cf-Connecting-Ip",
		"X-Client-Ip",
		"X-Forwarded-For",
		"X-Real-Ip",
	}

	for _, header := range headers {
		val := r.Header.Get(header)
		if val != "" {
			fmt.Fprintf(w, "%s: %v\n", header, val)
		}
	}
}

// Request dumps the HTTP request details.
func (h Handler) Request(w http.ResponseWriter, r *http.Request) {
	logger := Logger(r.Context())
	logger.Debug("Headers Handler")

	// show request values
	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("Error:\n%v\n", err),
			http.StatusInternalServerError,
		)
		logger.Error("DumpRequest", "err", err)
		return
	}

	fmt.Fprintln(w, string(b))
}
