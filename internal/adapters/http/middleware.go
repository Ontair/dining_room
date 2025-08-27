package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Middleware is func type that allows for
// chaining middleware
type Middleware func(http.HandlerFunc) http.HandlerFunc

type RequestCounter struct {
	mu        sync.RWMutex
	getCount  atomic.Uint64
	postCount atomic.Uint64
	getFile   string
	postFile  string
}

func NewRequestCounter(getFile, postFile string) *RequestCounter {
	return &RequestCounter{
		getFile:  getFile,
		postFile: postFile,
	}
}

func (rc *RequestCounter) IncrementGet() {
	rc.getCount.Add(1)
}

func (rc *RequestCounter) IncrementPost() {
	rc.postCount.Add(1)
}

func (rc *RequestCounter) WriteGetCount() error {
	count := rc.getCount.Load()

	content := fmt.Sprintf("GET Requests: %d\nLast Updated: %s\n",
		count, time.Now().UTC())
	
	rc.mu.Lock()
    defer rc.mu.Unlock()
	return os.WriteFile(rc.getFile, []byte(content), 0644)
}

func (rc *RequestCounter) WritePostCount() error {
	count := rc.postCount.Load()

	content := fmt.Sprintf("POST Requests: %d\nLast Updated: %s\n",
		count, time.Now().UTC())

	rc.mu.Lock()
    defer rc.mu.Unlock()
	return os.WriteFile(rc.postFile, []byte(content), 0644)
}

// CompileMiddleware takes the base http.HandlerFunc h
// and wraps it around the given list of Middleware m
func CompileMiddleware(h http.HandlerFunc, middlewareFunc []Middleware) http.HandlerFunc {
	if len(middlewareFunc) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(middlewareFunc) - 1; i >= 0; i-- {
		wrapped = middlewareFunc[i](wrapped)
	}

	return wrapped
}

func CreateCountRequestMiddleware(counter *RequestCounter) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				counter.IncrementGet()
				slog.Info("GET request counted", "total_get", counter.getCount.Load())
			case "POST":
				counter.IncrementPost()
				slog.Info("POST request counted", "total_post", counter.postCount.Load())
			}
			h.ServeHTTP(w, r)
		})
	}
}

func CreateWriteTxtMiddleware(counter *RequestCounter) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)

			switch r.Method {
			case "GET":
				if err := counter.WriteGetCount(); err != nil {
					slog.Error("Failed to write GET count to file", "error", err)
				}
			case "POST":
				if err := counter.WritePostCount(); err != nil {
					slog.Error("Failed to write POST count to file", "error", err)
				}
			}
		})
	}
}
