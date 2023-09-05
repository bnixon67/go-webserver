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
	"sort"
)

// HeadersPageName is the name of the HTTP template to execute.
const HeadersPageName = "headers.html"

// HeaderInfo contains individual header details.
type HeaderInfo struct {
	Key   string
	Value []string
}

// HeadersPageData holds the data passed to the HTML template.
type HeadersPageData struct {
	Title   string       // Title of the page.
	Headers []HeaderInfo // Sorted list of the request headers.
}

// NewHeaderInfo initializes and returns a sorted array of HeaderInfo from httpHeader.
func NewHeaderInfo(httpHeader http.Header) []HeaderInfo {
	headerList := make([]HeaderInfo, 0, len(httpHeader))
	for key, values := range httpHeader {
		headerList = append(headerList, HeaderInfo{Key: key, Value: values})
	}

	// Sort headers based on their key names
	sort.Slice(headerList, func(i, j int) bool { return headerList[i].Key < headerList[j].Key })

	return headerList
}

// HeadersHandler prints the headers of the request in sorted order.
func (h *Handler) HeadersHandler(w http.ResponseWriter, r *http.Request) {
	// get the logger from the context, which include request information
	logger := Logger(r.Context())
	logger.Info("HeadersHandler", "header", r.Header)

	// check for valid methods
	if !ValidMethod(w, r, http.MethodGet) {
		logger.Error("invalid method")
		return
	}

	sortedHeaders := NewHeaderInfo(r.Header)

	data := HeadersPageData{
		Title:   "Request Headers",
		Headers: sortedHeaders,
	}

	err := RenderTemplate(h.Tmpl, w, HeadersPageName, data)
	if err != nil {
		logger.Error("failed to RenderTemplate", "err", err)
		return
	}
}
