package main

import (
	"encoding/json"
	"errors"
	"net/http"
)


func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
  w.Header().Set("Content-Type", "application/json")
  data, err := json.Marshal(payload)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("something went wrong"))
    return
  }
  w.WriteHeader(status)
  w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
  respondWithJSON(w, code, errors.New(msg))
}
