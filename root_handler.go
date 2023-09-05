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
	"net/http"
)

// RootPageName is the name of the HTTP template to execute.
const RootPageName = "root.html"

// RootPageData contains any data passed to the HTML template.
type RootPageData struct {
	Title string // Title of the page.
}

// RootHandler handles the root ("/") route.
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	// get the logger from the context, which include request information
	logger := Logger(r.Context())

	// check for valid methods
	if !ValidMethod(w, r, http.MethodGet) {
		logger.Error("invalid method")
		return
	}

	// check for valid URL path
	if r.URL.Path != "/" {
		logger.Error("invalid path")
		http.NotFound(w, r)
		return
	}

	data := RootPageData{
		Title: h.AppName,
	}

	err := RenderTemplate(h.Tmpl, w, RootPageName, data)
	if err != nil {
		logger.Error("unable to RenderTemplate", "err", err)
		return
	}
}
