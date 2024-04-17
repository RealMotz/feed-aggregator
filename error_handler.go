package main

import (
	"encoding/json"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusInternalServerError)
  response := map[string]string{
    "error":"Internal Server Error",
  }
  data, err := json.Marshal(response)
  if err != nil {
    w.Write([]byte("error marshaling data"))
    return
  }
  w.Write(data)
}
