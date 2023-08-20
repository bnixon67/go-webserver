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

// RealIP returns the real IP address from the request headers or RemoteAddr.
func RealIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP")
	if realIP == "" {
		realIP = r.RemoteAddr
	}

	return realIP
}

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

		// add new Logger to context
		ctx := context.WithValue(r.Context(), loggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger returns the custom logger if present, otherwise default logger.
func Logger(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
