package main

import "net/http"

func defaultHandler(status int, headers map[string]string, body []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range headers {
			w.Header().Add(k, v)
		}
		w.WriteHeader(status)
		w.Write(body)
	})
}
