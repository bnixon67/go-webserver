package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"sort"
)

// Handler is responsible for handling HTTP requests.
type Handler struct{}

// Hello responds with a simple "hello" message.
func (h Handler) Hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v handler", r.URL.Path)

	fmt.Fprintln(w, "hello")
}

// Headers prints the headers of the request in sorted order.
func (h Handler) Headers(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v handler", r.URL.Path)

	// get keys
	keys := make([]string, 0, len(r.Header))
	for key := range r.Header {
		keys = append(keys, key)
	}

	// sort keys
	sort.Strings(keys)

	// print key and values
	for _, key := range keys {
		value := r.Header[key]
		fmt.Fprintf(w, "%v: %v\n", key, value)
	}
}

// IP responds with the remote IP address.
// This may not be the actual remote IP if a proxy, load balancer,
// or similar is used to route the request.
func (h Handler) IP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v handler", r.URL.Path)

	fmt.Fprintf(w, "RemoteAddr: %v\n", r.RemoteAddr)
}

// Request dumps the HTTP request details.
func (h Handler) Request(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v handler", r.URL.Path)

	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("Error:\n%v\n", err),
			http.StatusInternalServerError,
		)
		log.Print(err)
		return
	}

	fmt.Fprintln(w, string(b))
}

func main() {
	handler := Handler{}

	http.HandleFunc("/hello", handler.Hello)
	http.HandleFunc("/headers", handler.Headers)
	http.HandleFunc("/ip", handler.IP)
	http.HandleFunc("/request", handler.Request)

	addr := ":8080"
	log.Println("listen and serve on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
