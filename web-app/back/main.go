package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

  m "github.com/Seifbarouni/private-git/web-app/back/middlewares"

)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func InitHandlerPost(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// retrun the user as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	logger := slog.Default()

	mux := http.NewServeMux()

	mux.Handle("POST /",m.AuthorizationMiddleware(http.HandlerFunc(InitHandlerPost)))

	logger.Info("Starting server on port 8080")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		logger.Error("Failed to start server", "error", err)
	}

}
