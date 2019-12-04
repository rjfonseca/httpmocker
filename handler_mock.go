package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func mockHandler(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var resp response
		if err = json.NewDecoder(f).Decode(&resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(resp.StatusCode)
		for k, v := range resp.Header {
			for _, v := range v {
				w.Header().Add(k, v)
			}
		}
		w.Write(resp.Body)
	})
}
