package main

import (
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
		slog.Info("request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
		)

		handler.ServeHTTP(w, r)
	})
}

// Root handles the root ("/") route.
func (h Handler) Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/hello", http.StatusMovedPermanently)
		return
	}

	http.NotFound(w, r)
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

// IP responds with the remote IP address.
// Note: This may not be the actual remote IP if a proxy, load balancer,
// or similar is used to route the request.
func (h Handler) IP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RemoteAddr: %v\n", r.RemoteAddr)

	headers := []string{"X-Forwarded-For", "Cf-Connecting-Ip"}

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

	slog.Info("listen and serve", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		slog.Error("ListenAndServe", "err", err)
		os.Exit(1)
	}
}
