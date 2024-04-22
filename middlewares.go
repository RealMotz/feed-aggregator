package main

import (
	"net/http"
	"strings"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apikey := strings.Split(r.Header.Get("Authorization"), " ")
		if len(apikey) < 2 || apikey[0] != "ApiKey" {
			respondWithError(w, http.StatusUnauthorized, "unauthorized request")
			return
		}

		next.ServeHTTP(w, r)
	})
}
