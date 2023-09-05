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
)

// BuildHandler responds with the executable modification date and time.
func (h *Handler) BuildHandler(w http.ResponseWriter, r *http.Request) {
	logger := Logger(r.Context())

	if !ValidMethod(w, r, http.MethodGet) {
		logger.Error("invalid method")
		return
	}

	// try and force client not to cache content
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// get executable date/time
	dt, err := ExecutableDateTime()
	if err != nil {
		http.Error(w, MsgTemplateError, http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, dt.Format("2006-01-02 15:04:05"))
}
