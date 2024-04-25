package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RealMotz/feed-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couln't decode body")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    params.FeedId,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating record")
		return
	}

	respondWithJSON(w, http.StatusCreated, dbFeedFollowToFeedFollow(feedFollow))
}

func (cfg *apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid feed id")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     id,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete feed")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (cfg *apiConfig) getFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	follows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get feed follows")
		return
	}
	respondWithJSON(w, http.StatusOK, dbFeedFollowsToFeedFollows(follows))
}
