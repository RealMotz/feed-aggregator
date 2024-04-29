package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/RealMotz/feed-aggregator/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, maxFeedsToFetch int, intervalInSec int) {
	ticker := time.NewTicker(time.Duration(intervalInSec) * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			ProcessOldestFeedsFromDB(db, maxFeedsToFetch)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func ProcessOldestFeedsFromDB(db *database.Queries, maxFeedsToFetch int) {
	ctx := context.Background()
	feeds, err := db.GetNextFeedsToFetch(ctx, int32(maxFeedsToFetch))
	if err != nil {
		fmt.Println("error fetching feed")
		return
	}

	log.Printf("Found %v feeds to fetch!", len(feeds))

	group := &sync.WaitGroup{}

	for _, feed := range feeds {
		group.Add(1)
		go scrapeFeed(group, db, feed)
	}

	group.Wait()
}

func scrapeFeed(group *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer group.Done()
	items, err := fetchDataFromFeed(feed.Url)
	if err != nil {
		log.Printf("couldn't fetch feed %s", feed.Name)
		return
	}

	for _, item := range items {
		createPost(feed.ID, item, db)
	}

	_, err = db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed %s fetched %v", feed.Name, err)
		return
	}
}

func createPost(feedId uuid.UUID, item item, db *database.Queries) {
	var parsedPublicationDate time.Time
	rfc1123z, err := time.Parse(time.RFC1123Z, item.PublicationDate)
	if err != nil {
		rfc1123, err := time.Parse(time.RFC1123, item.PublicationDate)
		if err != nil {
			log.Printf("error parsing publication date: %v", err)
			return
		}
		parsedPublicationDate = rfc1123

	} else {
		parsedPublicationDate = rfc1123z
	}

	_, err = db.CreatePost(context.Background(), database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       item.Title,
		Url:         item.Link,
		Description: item.Description,
		PublisedAt:  parsedPublicationDate,
		FeedID:      feedId,
	})

	if err != nil {
		return
	}
}

func fetchDataFromFeed(feedURL string) ([]item, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		fmt.Println("error fetching feed")
		return nil, err
	}

	defer resp.Body.Close()

	data := xmlData{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Channel.ItemList, nil
}

type item struct {
	Title           string         `xml:"title"`
	Link            string         `xml:"link"`
	PublicationDate string         `xml:"pubDate"`
	Description     sql.NullString `xml:"description"`
}

type xmlEntry struct {
	ItemList []item `xml:"item"`
}

type xmlData struct {
	Channel xmlEntry `xml:"channel"`
}
