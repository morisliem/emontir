package handler

import (
	"context"
	"e-montir/pkg/jwt"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ContextKey string

const (
	tokenKey = ContextKey("token")
)

func DecodeJSON(r *http.Request, data interface{}) error {
	if r.Body == nil {
		return nil
	}
	defer r.Body.Close()

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		return fmt.Errorf("content-type not found")
	}

	return json.NewDecoder(r.Body).Decode(data)
}

func GetTokenClaim(ctx context.Context) *jwt.Claim {
	return ctx.Value(tokenKey).(*jwt.Claim)
}
