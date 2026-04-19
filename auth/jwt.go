package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken   = errors.New("invalid or expired token")
	ErrMissingKeyfunc = errors.New("JWKS keyfunc is not initialized")
)

// Validator handles validating JWT tokens using an OIDC JWKS endpoint
type Validator struct {
	jwks *keyfunc.JWKS
}

// NewValidator initializes a new Validator by fetching the JWKS from the given endpoint
// Typically the URL resembles: https://keycloak.example.com/realms/halooid/protocol/openid-connect/certs
func NewValidator(ctx context.Context, jwksURL string, refreshInterval time.Duration) (*Validator, error) {
	options := keyfunc.Options{
		Ctx:                 ctx,
		RefreshInterval:     refreshInterval,
		RefreshErrorHandler: func(err error) {
			// Log error via provided logger (simplifying for go-shared)
			fmt.Printf("failed to refresh jwks: %s\n", err)
		},
	}

	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS from %s: %w", jwksURL, err)
	}

	return &Validator{
		jwks: jwks,
	}, nil
}

// Validate takes a raw JWT token and returns the parsed Claims. It verifies the signature using RS256.
func (v *Validator) Validate(rawToken string) (*Claims, error) {
	if v.jwks == nil {
		return nil, ErrMissingKeyfunc
	}

	token, err := jwt.ParseWithClaims(rawToken, &Claims{}, v.jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ExtractBearerToken parses the Bearer string and returns the raw token
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}
