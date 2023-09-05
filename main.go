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
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"
)

const (
	ExitServer   = 1 // ExitServer indicates a server error.
	ExitUsage        // ExitUsage indicates a usage error.
	ExitLog          // ExitLog indicates a log error.
	ExitTemplate     // ExitTemplate indicates a template error.
)

// ServerConfig holds configuration options for the HTTP server.
type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// createServer creates an HTTP server with the specified config and handler.
func createServer(config ServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         config.Addr,
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}
}

// runServer starts the HTTP server and handles graceful shutdown.
func runServer(ctx context.Context, srv *http.Server) {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		slog.Error("failed to listen", "err", err)
		os.Exit(ExitServer)
	}

	go func() {
		err := srv.Serve(ln)
		if err != nil && err != http.ErrServerClosed {
			slog.Error("failed to serve", "err", err)
			os.Exit(ExitServer)
		}
	}()

	slog.Info("started server", slog.String("addr", ln.Addr().String()))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		signal.Stop(sigChan)

		slog.Info("shutting down server", "signal", sig)

		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err := srv.Shutdown(timeoutCtx)
		if err != nil {
			slog.Error("server shutdown error", "err", err)
		}

		slog.Info("server shutdown")
	case <-ctx.Done():
	}
}

func main() {
	// define command-line flags
	addrFlag := flag.String("addr", "localhost:8080", "[host]:port")
	logFileFlag := flag.String("logfile", "", "log file")
	logLevelFlag := flag.String("loglevel", "Info", "log level")
	logTypeFlag := flag.String("logtype", "json", "log type (json|text)")
	logAddSource := flag.Bool("logsource", false, "log source code position")

	// parse command-line flags
	flag.Parse()

	// get slog.Level from flag
	logLevel, err := LogLevel(*logLevelFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fmt.Fprintf(os.Stderr, "valid loglevels: %s\n", LogLevels())
		flag.Usage()
		os.Exit(ExitUsage)
	}

	// validate logtype
	if !slices.Contains(validLogTypes, *logTypeFlag) {
		fmt.Fprintf(os.Stderr, "invalid logtype: %v\n", *logTypeFlag)
		fmt.Fprintf(os.Stderr, "valid logtypes: %s\n", strings.Join(validLogTypes, ", "))
		flag.Usage()
		os.Exit(ExitUsage)
	}

	// check for additional command-line arguments
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(ExitUsage)
	}

	// initialize logging
	err = InitLog(*logFileFlag, *logTypeFlag, logLevel, *logAddSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(ExitLog)
	}

	// initialize templates
	tmpl, err := InitTemplates("html/*.html")
	if err != nil {
		slog.Error("failed to InitTemplates", "err", err)
		os.Exit(ExitTemplate)
	}

	h := NewHandler("Go Web Server", tmpl)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.RootHandler)
	mux.HandleFunc("/hello", h.HelloHandler)
	mux.HandleFunc("/hellohtml", h.HelloHTMLHandler)
	mux.HandleFunc("/headers", h.HeadersHandler)
	mux.HandleFunc("/remote", h.RemoteHandler)
	mux.HandleFunc("/request", h.RequestHandler)
	mux.HandleFunc("/build", h.BuildHandler)

	ctx := context.Background()
	serverConfig := ServerConfig{
		Addr:         *addrFlag,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	srv := createServer(serverConfig, h.AddRequestID(h.LogRequest(mux)))
	runServer(ctx, srv)
}
