package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func withSleep(sleep, jitter time.Duration, next http.Handler) http.Handler {
	rand.Seed(time.Now().UnixNano())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sleepTime := sleep + time.Duration(rand.Intn(int(jitter)))
		log.Println("Sleeping", sleepTime)
		time.Sleep(sleepTime)
		next.ServeHTTP(w, r)

	})
}
