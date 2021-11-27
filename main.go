package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
)

func hello(w http.ResponseWriter, req *http.Request) {
	log.Println("hello handler")

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	log.Println("headers handler")

	// Get keys from map.
	keys := []string{}
	for key, _ := range req.Header {
		keys = append(keys, key)
	}

	// Sort string keys.
	sort.Strings(keys)

	// Loop over sorted key-value pairs.
	for i := range keys {
		key := keys[i]
		value := req.Header[key]
		fmt.Fprintf(w, "%v: %v\n", key, value)
	}
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	addr := ":8080"
	log.Println("listen and serve on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
