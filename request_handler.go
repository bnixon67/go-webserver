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
	"fmt"
	"net/http"
	"net/http/httputil"
)

// RequestHandler dumps the HTTP request details.
func (h *Handler) RequestHandler(w http.ResponseWriter, r *http.Request) {
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
