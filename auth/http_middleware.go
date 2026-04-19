package auth

import (
	"context"
	"net/http"
)

func Middleware(v *Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := ExtractBearerToken(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			claims, err := v.Validate(tokenStr)
			if err != nil {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			ctx := WithClaims(r.Context(), claims)
			r = r.WithContext(ctx)

			// Propagate gateway headers conventionally if needed downstream
			r.Header.Set("x-user-id", claims.Subject)
			r.Header.Set("x-user-email", claims.Email)
			r.Header.Set("x-tenant-id", claims.TenantID)

			next.ServeHTTP(w, r)
		})
	}
}
