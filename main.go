package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"

	"golang.org/x/exp/slog"
)

// Handler is responsible for handling HTTP requests.
type Handler struct{}

// logger is HTTP middleware to log the request.
func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get real IP address if using Cloudflare or similar service
		var ip string
		ip = r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = r.RemoteAddr
		}

		slog.Info("request",
			slog.String("ip", ip),
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
		)

		handler.ServeHTTP(w, r)
	})
}

//go:embed testdata/hello.html
var helloHTML string

// Root handles the root ("/") route.
func (h Handler) Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprint(w, helloHTML)
}

// Hello responds with a simple "hello" message.
func (h Handler) Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Fprintln(w, "hello")
}

// Headers prints the headers of the request in sorted order.
func (h Handler) Headers(w http.ResponseWriter, r *http.Request) {
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

// IP responds with the remote IP and common headers for the actual IP.
// Note: RemoteAddr may not be the actual remote IP if a proxy, load balancer,
// or similar is used to route the request.
func (h Handler) IP(w http.ResponseWriter, r *http.Request) {
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
	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("Error:\n%v\n", err),
			http.StatusInternalServerError,
		)
		slog.Error("DumpRequest", "err", err)
		return
	}

	fmt.Fprintln(w, string(b))
}

func main() {
	// Define a flag for the address
	addr := flag.String("addr", ":8080", "address (host:port) to listen on")
	flag.Parse()

	var handler Handler

	mux := http.NewServeMux()

	mux.Handle("/", logger(http.HandlerFunc(handler.Root)))
	mux.Handle("/hello", logger(http.HandlerFunc(handler.Hello)))
	mux.Handle("/headers", logger(http.HandlerFunc(handler.Headers)))
	mux.Handle("/ip", logger(http.HandlerFunc(handler.IP)))
	mux.Handle("/request", logger(http.HandlerFunc(handler.Request)))

	slog.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("server closed")
	} else if err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}
