package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/RealMotz/feed-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	dbURl := os.Getenv("CONN")
	db, err := sql.Open("postgres", dbURl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	cfg := apiConfig{
		DB: dbQueries,
	}

	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: corsMux,
	}

	go createInterval(10, cfg.ProcessOldestNFeedsFromDB)

	mux.HandleFunc("GET /v1/readiness", readinessHandler)
	mux.HandleFunc("GET /v1/err", errorHandler)

	mux.HandleFunc("POST /v1/users", cfg.createUser)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.getUser))

	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.createFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.getFeeds)

	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.createFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{id}", cfg.middlewareAuth(cfg.deleteFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.getFeedFollow))

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func createInterval(seconds int, fn func(int32)) {
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			fn(10)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (cfg *apiConfig) ProcessOldestNFeedsFromDB(n int32) {
	ctx := context.Background()
	feeds, err := cfg.DB.GetNextFeedsToFetch(ctx, n)
	if err != nil {
		fmt.Println("error fetching feed")
		return
	}

	var group sync.WaitGroup

	for _, feed := range feeds {
		group.Add(1)
		go func(url string) {
			defer group.Done()
			fetchDataFromFeed(url)
		}(feed.Url)
	}

	group.Wait()
}

func fetchDataFromFeed(feedURL string) {
	resp, err := http.Get(feedURL)
	if err != nil {
		fmt.Println("error fetching feed")
		return
	}
	defer resp.Body.Close()

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

	data := xmlData{}
	decoder := xml.NewDecoder(resp.Body)
	decoder.Decode(&data)

	fmt.Printf("processing %s\n", feedURL)
	fmt.Println(len(data.Channel.ItemList))
}
