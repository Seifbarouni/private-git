package middlewares

import (
  "net/http"
  "fmt"
  "errors"
  "log/slog"
  "os"
  "context"

  "github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "authorization header not provided", http.StatusUnauthorized)
			return
		}

    token, err := verifyToken(tokenString[7:])
    if err != nil {
      http.Error(w, "invalid token", http.StatusUnauthorized)
      return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
      http.Error(w, "invalid token claims", http.StatusUnauthorized)
      return
    }

    slog.Info(fmt.Sprintf("User %s with email: %s is authorized", claims["name"], claims["email"]))

    r = r.WithContext(context.WithValue(r.Context(), "name", claims["name"]))
    r = r.WithContext(context.WithValue(r.Context(), "email", claims["email"]))
    
		next.ServeHTTP(w, r)
	})
}

func verifyToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    return token, nil
}
