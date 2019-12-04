package main

import "net/http"

func multiHandler(maxLoops int, defaultHandler http.Handler, handlers ...http.Handler) http.Handler {
	nextHandler := make(chan http.Handler)
	go func() {
		if len(handlers) > 0 {
			for loopCount := 0; maxLoops == 0 || loopCount < maxLoops; loopCount++ {
				for _, h := range handlers {
					nextHandler <- h
				}
			}
		}
		for {
			nextHandler <- defaultHandler
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := <-nextHandler
		h.ServeHTTP(w, r)
	})
}
