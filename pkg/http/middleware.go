package httpmiddleware

import (
    "net/http"
    
    "github.com/amitstephen-dev/goheapguard/pkg/heapguard"
)

func WithPooledHandler[T any](pool *heapguard.Pool[T], handler func(*T, http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        obj := pool.Get()
        defer pool.Put(obj)

        handler(obj, w, r)
    }
}