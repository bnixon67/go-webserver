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
)

// HelloHandler responds with a simple "hello" message in text format.
func (h *Handler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	logger := Logger(r.Context())
	logger.Debug("Hello Handler")

	if !ValidMethod(w, r, http.MethodGet) {
		logger.Error("invalid method")
		return
	}

	// try and force client not to cache content
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Fprintln(w, "hello")
}

//go:embed html/hello.html
var helloHTML string

// HelloHTMLHandler responds with a simple "hello" message in HTML format.
func (h *Handler) HelloHTMLHandler(w http.ResponseWriter, r *http.Request) {
	logger := Logger(r.Context())
	logger.Debug("Hello HTML Handler")

	if !ValidMethod(w, r, http.MethodGet) {
		logger.Error("invalid method")
		return
	}

	// try and force client not to cache content
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// Write the HTML content to the response
	fmt.Fprint(w, helloHTML)
}
