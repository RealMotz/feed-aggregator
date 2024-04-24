package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/RealMotz/feed-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot fetch feed data")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}

func (cfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string
		Url  string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couln't decode body")
		return
	}

	// TODO: Refactor this section
	apikey := strings.Split(r.Header.Get("Authorization"), " ")[1]

	user, err := cfg.DB.GetUserByApikey(r.Context(), apikey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, dbFeedToFeed(feed))
}
