package main

import (
	"net/http"
	"strconv"

	"github.com/RealMotz/feed-aggregator/internal/database"
)

func (cfg *apiConfig) getPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	limit, err := getLimit(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid query parameter")
		return
	}
	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching posts")
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func getLimit(r *http.Request) (int, error) {
	query := r.URL.Query().Get("limit")
	if query == "" {
		return 10, nil
	}

	limit, err := strconv.Atoi(query)
	if err != nil {
		return 0, err
	}

	return limit, nil
}
