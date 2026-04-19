package auth

import "github.com/golang-jwt/jwt/v5"

// Claims represents the structure of the JWT provided by Keycloak
type Claims struct {
	jwt.RegisteredClaims
	Email             string   `json:"email"`
	PreferredUsername string   `json:"preferred_username"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	TenantID          string   `json:"tenant_id,omitempty"` // Custom Keycloak attribute if configured
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}
