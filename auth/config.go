// Package auth contain wrapper for keycloak client
package auth

type Config struct {
	Keycloak KeycloakConfig `json:"keycloak" yaml:"keycloak" envPrefix:"KEYCLOAK_"`
}

type KeycloakConfig struct {
	BaseURL string      `json:"base_url" yaml:"base_url" env:"BASE_URL"`
	Realm   string      `json:"realm" yaml:"realm" env:"REALM"`
	Admin   Credentials `json:"admin" yaml:"admin" envPrefix:"ADMIN_"`
}

type Credentials struct {
	Username     string `json:"username" yaml:"username" env:"USERNAME"`
	Password     string `json:"password" yaml:"password" env:"PASSWORD"`
	ClientID     string `json:"client_id" yaml:"client_id" env:"CLIENT_ID"`
	ClientSecret string `json:"client_secret" yaml:"client_secret" env:"CLIENT_SECRET"`
}
