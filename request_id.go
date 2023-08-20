package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync/atomic"
)

type ctxKey int

const requestIDKey ctxKey = iota

var reqIDPrefix string

// randomString generates a random string of the specified length, composed of
// uppercase letters, lowercase letters, and digits.
func randomString(length int) (string, error) {
	const (
		upper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lower  = "abcdefghijklmnopqrstuvwxyz"
		digits = "0123456789"
		chars  = upper + lower + digits
	)

	// check for valid length
	if length <= 0 {
		return "", errors.New("invalid length")
	}

	result := make([]byte, length)

	for i := 0; i < length; i++ {
		// generate random index
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[idx.Int64()]
	}

	return string(result), nil
}

// init generates a unique request id prefix at program start
func init() {
	reqIDPrefix, _ = randomString(6)
}

// geneateRequestID returns the combination of request id prefix and counter.
func generateRequestID(counter uint32) string {
	return fmt.Sprintf("%s%010d", reqIDPrefix, counter)
}

// AddRequestID is middleware that adds, for each request, a unique
// incrementing request ID to the context and headers.
func (h Handler) AddRequestID(next http.Handler) http.Handler {
	var counter uint32

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID(atomic.AddUint32(&counter, 1))
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext returns request ID from ctx if present, otherwise "".
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	reqID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return reqID
}
