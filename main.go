package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
)

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("hello handler", r.Method, r.URL)

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	fmt.Fprintln(w, "hello")
}

func root(w http.ResponseWriter, r *http.Request) {
	log.Println("root handler", r.Method, r.URL)

	if r.URL.Path == "/" {
		http.Redirect(w, r, "/hello", http.StatusMovedPermanently)
		return
	} else {
		http.NotFound(w, r)
		return
	}
}

func headers(w http.ResponseWriter, r *http.Request) {
	log.Println("headers handler", r.Method, r.URL)

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// get keys from map
	keys := []string{}
	for key, _ := range r.Header {
		keys = append(keys, key)
	}

	// sort keys
	sort.Strings(keys)

	// print key and values in sorted order
	for i := range keys {
		key := keys[i]
		value := r.Header[key]
		fmt.Fprintf(w, "%v: %v\n", key, value)
	}
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	addr := ":8080"
	log.Println("listen and serve on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
