package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/RealMotz/feed-aggregator/internal/database"
)

func startScraping(db *database.Queries, maxFeeds int, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			ProcessOldestFeedsFromDB(db, maxFeeds)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func ProcessOldestFeedsFromDB(db *database.Queries, maxFeeds int) {
	ctx := context.Background()
	feeds, err := db.GetNextFeedsToFetch(ctx, int32(maxFeeds))
	if err != nil {
		fmt.Println("error fetching feed")
		return
	}

	group := &sync.WaitGroup{}

	for _, feed := range feeds {
		group.Add(1)
		go scrapeFeed(group, db, feed)
	}

	group.Wait()
}

func scrapeFeed(group *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer group.Done()
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed %s fetched %v", feed.Name, err)
		return
	}

	items, err := fetchDataFromFeed(feed.Url)
	if err != nil {
		log.Printf("couldn't fetch feed %s", feed.Name)
		return
	}

	fmt.Printf("processing %s\n", feed.Url)
	fmt.Println(len(items))
}

func fetchDataFromFeed(feedURL string) ([]itemList, error) {
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

type itemList struct {
	Title           string `xml:"title"`
	Link            string `xml:"link"`
	PublicationDate string `xml:"pubDate"`
	Description     string `xml:"description"`
}

type xmlEntry struct {
	ItemList []itemList `xml:"item"`
}

type xmlData struct {
	Channel xmlEntry `xml:"channel"`
}
