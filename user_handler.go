package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/RealMotz/feed-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couln't decode body")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "can't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, dbUserToUser(user))
}

func (cfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request) {
	apikey := strings.Split(r.Header.Get("Authorization"), " ")[1]

	user, err := cfg.DB.GetUserByApikey(r.Context(), apikey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user not found")
		return
	}

	respondWithJSON(w, http.StatusOK, dbUserToUser(user))
}
