package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func registerRequest(outputdir string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, ok := r.Context().Value(requestNumberKey).(uint64)
		if !ok {
			log.Fatalln("missing request number")
		}
		request := newRequestFromHTTP(r)
		b, err := json.Marshal(request)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(filepath.Join(outputdir, fmt.Sprintf("%d.request.json", n)), b, 0666)
		if err != nil {
			panic(err)
		}
		next.ServeHTTP(w, r)
	})
}
