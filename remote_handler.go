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

// RemoteHandler responds with the RemoteHandler and common headers for the actual RemoteHandler.
// Note: RemoteHandler may not be valid if a proxy, load balancer, or similar is used to route the request.
func (h *Handler) RemoteHandler(w http.ResponseWriter, r *http.Request) {
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
