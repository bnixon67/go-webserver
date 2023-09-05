package main

import (
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

// ValidMethod checks if the given HTTP request method is allowed based on
// the provided list of allowed methods. It returns true if the method is
// allowed, and false otherwise. If the method is not allowed or is OPTIONS, the
// function updates the response writer appropriately and returns false.  The
// calling handler should return without further processing.
func ValidMethod(w http.ResponseWriter, r *http.Request, allowed ...string) bool {
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

// ExecutableDateTime returns the modification date/time of the executable file.
func ExecutableDateTime() (time.Time, error) {
	var time time.Time

	// Get the path of the executable
	executablePath, err := os.Executable()
	if err != nil {
		return time, err
	}

	// Get file information
	fileInfo, err := os.Stat(executablePath)
	if err != nil {
		return time, err
	}

	return fileInfo.ModTime(), nil
}

// RealIP returns the real IP address from the request headers or RemoteAddr.
func RealIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP")
	if realIP == "" {
		realIP = r.RemoteAddr
	}

	return realIP
}
