package main

import (
	"time"

	"github.com/RealMotz/feed-aggregator/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name"`
	Url         string     `json:"url"`
	UserID      uuid.UUID  `json:"user_id"`
	LastFetchAt *time.Time `json:"last_fetch_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Apikey    string    `json:"api_key"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	FeedId    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func dbUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		Apikey:    dbUser.Apikey,
	}
}

func dbFeedToFeed(dbFeed database.Feed) Feed {
	var lastFetch *time.Time = nil
	if dbFeed.LastFetchAt.Valid {
		lastFetch = &dbFeed.LastFetchAt.Time
	}
	return Feed{
		ID:          dbFeed.ID,
		CreatedAt:   dbFeed.CreatedAt,
		UpdatedAt:   dbFeed.UpdatedAt,
		Name:        dbFeed.Name,
		Url:         dbFeed.Url,
		UserID:      dbFeed.UserID,
		LastFetchAt: lastFetch,
	}
}

func dbFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		FeedId:    dbFeedFollow.FeedID,
		UserID:    dbFeedFollow.UserID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
	}
}

func dbFeedsToFeeds(feeds []database.Feed) []Feed {
	result := make([]Feed, len(feeds))
	for i, feed := range feeds {
		result[i] = dbFeedToFeed(feed)
	}
	return result
}

func dbFeedFollowsToFeedFollows(follows []database.FeedFollow) []FeedFollow {
	result := make([]FeedFollow, len(follows))
	for i, follow := range follows {
		result[i] = dbFeedFollowToFeedFollow(follow)
	}
	return result
}
