package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Session struct {
	jwt.StandardClaims
	SessionState      string      `json:"session_state"`
	Acr               string      `json:"acr"`
	Aud               interface{} `json:"aud"`
	Scope             string      `json:"scope"`
	EmailVerified     bool        `json:"email_verified"`
	Name              string      `json:"name"`
	PreferredUsername string      `json:"preferred_username"`
	GivenName         string      `json:"given_name"`
	FamilyName        string      `json:"family_name"`
	Email             string      `json:"email"`
	Groups            []string    `json:"groups"`
}
