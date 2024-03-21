package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for an Authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Authorization header not provided", http.StatusUnauthorized)
			return
		}

		// TODO: Validate the token. This is just a placeholder.
		if token != "Bearer hello" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// If we got this far, the token is valid and we can call the next handler
		next.ServeHTTP(w, r)
	})
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

	mux.Handle("POST /", AuthorizationMiddleware(http.HandlerFunc(InitHandlerPost)))

	logger.Info("Starting server on port 8080")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		logger.Error("Failed to start server", "error", err)
	}

}
