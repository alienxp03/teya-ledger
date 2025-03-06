package api

import (
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// resolve userID
		// ideally we should query from the database to get the userID
		userID := strings.ReplaceAll(authHeader, "USER_TOKEN", "USER_ID")

		ctx := context.WithValue(r.Context(), HeaderUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
