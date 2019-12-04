package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type response struct {
	request
	StatusCode int
}

func registerResponse(outputdir string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, ok := r.Context().Value(requestNumberKey).(uint64)
		if !ok {
			log.Fatalln("missing request number")
		}

		copy := &copyResponseWriter{
			Writer: w,
			Buffer: new(bytes.Buffer),
		}

		next.ServeHTTP(copy, r)

		b, err := json.Marshal(response{
			request: request{
				URL:    r.URL.String(),
				Body:   copy.Buffer.Bytes(),
				Header: w.Header(),
			},
			StatusCode: copy.StatusCode,
		})
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(filepath.Join(outputdir, fmt.Sprintf("%d.response.json", n)), b, 0666)
		if err != nil {
			panic(err)
		}
	})
}

type copyResponseWriter struct {
	Writer     http.ResponseWriter
	Buffer     *bytes.Buffer
	StatusCode int
}

func (w *copyResponseWriter) Header() http.Header {
	return w.Writer.Header()
}

func (w *copyResponseWriter) Write(body []byte) (int, error) {
	w.Buffer.Write(body)
	return w.Writer.Write(body)
}

func (w *copyResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.Writer.WriteHeader(statusCode)
}
