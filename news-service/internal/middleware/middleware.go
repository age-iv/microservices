package middleware

import (
    "context"
    "log"
    "net/http"
    "time"
    "math/rand"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rid := r.URL.Query().Get("request_id")
        if rid == "" {
            rid = generateID(8)
        }
        ctx := context.WithValue(r.Context(), RequestIDKey, rid)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(wrapped, r)
        log.Printf("[%s] %s %s %d %v %s",
            r.RemoteAddr, r.Method, r.URL.Path,
            wrapped.statusCode, time.Since(start),
            r.Context().Value(RequestIDKey))
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func generateID(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
