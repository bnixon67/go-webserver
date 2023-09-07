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
	"context"
	"log/slog"
	"net/http"
)

// LoggerKey is used as a context key for the custom logger.
type LoggerKey int

const loggerKey LoggerKey = iota

// LogRequest middleware logs incoming HTTP requests.
// It adds a Logger to the request context that can be used by child handlers to include request information.
// If the header X-Real-IP exists, it is used instead of RemoteAddr.
func (h Handler) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.With(slog.Group("request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("ip", RealIP(r)),
			slog.String("requestID", RequestIDFromContext(r.Context())),
		))
		logger.Info("LogRequest")

		// add new logger to context
		ctx := context.WithValue(r.Context(), loggerKey, logger)

		// server the request with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger returns the custom logger if present, otherwise default logger.
func Logger(ctx context.Context) *slog.Logger {
	// if context is nil, return the default logger.
	if ctx == nil {
		return slog.Default()
	}

	// attempt to retrieve the logger from the context using the loggerKey
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)

	// if the logger is not present in the context or has an incorrect type, return the default logger.
	if !ok {
		return slog.Default()
	}

	// return the custom logger from the context
	return logger
}
