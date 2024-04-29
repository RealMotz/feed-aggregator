package main

import (
	"encoding/json"
	"net/http"
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

	respondWithJSON(w, http.StatusOK, dbFeedsToFeeds(feeds))
}

func (cfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
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

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating record")
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}{
		Feed:       dbFeedToFeed(feed),
		FeedFollow: dbFeedFollowToFeedFollow(feedFollow),
	})
}
