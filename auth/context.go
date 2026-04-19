package auth

import (
	"context"
	"errors"
)

type contextKey string

const (
	claimsContextKey contextKey = "authClaims"
)

var (
	ErrNoClaimsInContext = errors.New("no auth claims found in context")
)

// WithClaims returns a new context with the provided claims.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// GetUser retrieves Claims from the context, if present.
func GetUser(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsContextKey).(*Claims)
	return c, ok
}

// MustUser retrieves Claims from the context or panics. Use with caution.
func MustUser(ctx context.Context) *Claims {
	c, ok := GetUser(ctx)
	if !ok {
		panic(ErrNoClaimsInContext)
	}
	return c
}
