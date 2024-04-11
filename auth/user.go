package auth

import (
	"time"

	gocloak "github.com/Nerzal/gocloak/v6"
)

type User struct {
	ID               string    `json:"id"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
	Username         string    `json:"username"`
	Enabled          bool      `json:"enabled"`
	EmailVerified    bool      `json:"email_verified"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
}

func NewUserFromKeycloakModel(keycloakUser *gocloak.User) *User {
	if keycloakUser == nil {
		return nil
	}

	user := new(User)
	user.ID = *keycloakUser.ID
	if keycloakUser.CreatedTimestamp != nil {
		user.CreatedTimestamp = time.Unix(0, *keycloakUser.CreatedTimestamp*int64(time.Millisecond))
	}

	if keycloakUser.Username != nil {
		user.Username = *keycloakUser.Username
	}

	if keycloakUser.Enabled != nil {
		user.Enabled = *keycloakUser.Enabled
	}

	if keycloakUser.EmailVerified != nil {
		user.EmailVerified = *keycloakUser.EmailVerified
	}

	if keycloakUser.FirstName != nil {
		user.FirstName = *keycloakUser.FirstName
	}

	if keycloakUser.LastName != nil {
		user.LastName = *keycloakUser.LastName
	}

	if keycloakUser.Email != nil {
		user.Email = *keycloakUser.Email
	}

	return user
}
