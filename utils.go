package main

import (
	"net/http"
	"slices"
	"strings"
)

// ValidMethod checks if the given HTTP request method is allowed based on
// the provided list of allowed methods. It returns true if the method is
// allowed, and false otherwise. If the method is not allowed or is OPTIONS, the
// function updates the response writer appropriately and returns false.  The
// calling handler should return without further processing.
func ValidMethod(w http.ResponseWriter, r *http.Request, allowed ...string) bool {
	logger := Logger(r.Context())
	logger.Debug("ValidMethod", "allowed", allowed)

	// if method is in allowed list, then return
	if slices.Contains(allowed, r.Method) {
		return true
	}

	// add OPTIONS if method is not allowed
	allowed = append(allowed, http.MethodOptions)

	// set the "Allow" header to allowed methods
	w.Header().Set("Allow", strings.Join(allowed, ", "))

	// if method is OPTIONS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // no content returned
		return false
	}

	// method is not allowed and not OPTIONS
	txt := r.Method + " " + http.StatusText(http.StatusMethodNotAllowed)
	http.Error(w, txt, http.StatusMethodNotAllowed)
	return false
}
