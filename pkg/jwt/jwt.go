package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

type Claim struct {
	jwt.StandardClaims
	ID string `json:"id"`
}

func GenerateToken(id, key string, duration int) (token, expiredAt string, err error) {
	claim := Claim{
		ID: id,
	}
	now := time.Now().UTC()
	claim.IssuedAt = now.Unix()
	claim.ExpiresAt = now.Add(time.Minute * time.Duration(duration)).Unix()

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err = newToken.SignedString([]byte(key))
	if err != nil {
		log.Error().Err(fmt.Errorf("error when SignedToken: %w", err))
		return
	}

	expiredAt = now.Add(time.Minute * time.Duration(duration)).Format(time.RFC3339)
	return
}

func ParseTokenClaim(token, key string) (*Claim, error) {
	claim := new(Claim)
	tokenClaim, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		log.Error().Err(fmt.Errorf("error when ParsingToken: %w", err))
		return nil, err
	}

	if tokenClaim.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, fmt.Errorf("invalid signing method")
	}

	return claim, nil
}
