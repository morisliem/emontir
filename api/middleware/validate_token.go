package middleware

import (
	"context"
	"e-montir/api/handler"
	"e-montir/pkg/jwt"
	"net/http"
	"os"
	"strings"
)

const (
	tokenKey = handler.ContextKey("token")
)

func ValidateToken() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader == "" {
				handler.GenerateResponse(w, http.StatusUnauthorized, handler.UnauthorizedError)
				return
			}

			bearerToken := strings.Split(authorizationHeader, " ")
			token := bearerToken[1]

			if len(bearerToken) > 2 {
				handler.GenerateResponse(w, http.StatusUnauthorized, handler.UnauthorizedError)
				return
			}

			if len(bearerToken) != 2 && bearerToken[0] != "Bearer" {
				handler.GenerateResponse(w, http.StatusUnauthorized, handler.UnauthorizedError)
				return
			}

			claim, err := jwt.ParseTokenClaim(token, os.Getenv("ACCESS_KEY"))
			if err != nil {
				handler.GenerateResponse(w, http.StatusUnauthorized, handler.UnauthorizedError)
				return
			}

			ctx := context.WithValue(r.Context(), tokenKey, claim)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
