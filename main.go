package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

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

	go startScraping(dbQueries, 10, 10)

	mux.HandleFunc("GET /v1/readiness", readinessHandler)
	mux.HandleFunc("GET /v1/err", errorHandler)

	mux.HandleFunc("POST /v1/users", cfg.createUser)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.getUser))

	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.createFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.getFeeds)

	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.createFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{id}", cfg.middlewareAuth(cfg.deleteFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.getFeedFollow))

	mux.HandleFunc("GET /v1/feeds/fetch", cfg.middlewareAuth(cfg.fetchDataFromFeed2))

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (cfg *apiConfig) fetchDataFromFeed2(w http.ResponseWriter, r *http.Request, user database.User) {
	resp, err := http.Get("https://blog.boot.dev/index.xml")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching feed")
		fmt.Printf("Error GET: %v\n", err)
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

	respondWithJSON(w, http.StatusCreated, data.Channel.ItemList)
}
