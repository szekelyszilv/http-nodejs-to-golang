package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func CreateRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/golang", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "Hello from Golang",
		})
	})
	return mux
}
