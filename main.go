package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

  port := os.Getenv("PORT")
  mux := http.NewServeMux()
  corsMux := middlewareCors(mux)

  server := &http.Server{
    Addr: fmt.Sprintf(":%s", port),
    Handler: corsMux,
  }

  mux.HandleFunc("GET /v1/readiness", readinessHandler)
  mux.HandleFunc("GET /v1/err", errorHandler)

  err = server.ListenAndServe()
  if err != nil {
    log.Fatal(err)
  }
}
