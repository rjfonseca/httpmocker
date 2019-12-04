package main

import (
	"bytes"
	"net/http"
)

type request struct {
	Header map[string][]string
	URL    string
	Body   []byte
}

func newRequestFromHTTP(r *http.Request) request {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	/*
		header := map[string][]string
		for k, v := range r.Header {
			req.Header[k]=  v
		}
	*/
	return request{
		Header: r.Header,
		URL:    r.URL.String(),
		Body:   buf.Bytes(),
	}
}
