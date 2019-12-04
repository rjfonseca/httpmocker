package main

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"
)

func requestCounter(next http.Handler) http.Handler {
	var counter uint64
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := atomic.AddUint64(&counter, 1)
		log.Println("Processing request", c)
		ctx := context.WithValue(r.Context(), requestNumberKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
