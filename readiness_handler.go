package main

import (
	"encoding/json"
	"net/http"
)


func readinessHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  response := map[string]string{
    "status":"ok",
  }
  data, err := json.Marshal(response)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("error marshaling data"))
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(data)
}
