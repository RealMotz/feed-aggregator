package main

import (
	"net/http"
	"strings"

	"github.com/RealMotz/feed-aggregator/internal/database"
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

type authHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apikey := strings.Split(r.Header.Get("Authorization"), " ")
		if len(apikey) < 2 || apikey[0] != "ApiKey" {
			respondWithError(w, http.StatusUnauthorized, "unauthorized request")
			return
		}

		user, err := cfg.DB.GetUserByApikey(r.Context(), apikey[1])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "user not found")
			return
		}

		handler(w, r, user)
	})
}
