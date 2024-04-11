package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	gocloak "github.com/Nerzal/gocloak/v6"
	jwt "github.com/dgrijalva/jwt-go"
)

const TokenMinuteCorrectionFactor = 20

type AuthorizeFunc func(token string, realm string) error

type Client struct {
	Cfg         *Config
	PathRoleMap map[string][]string
	adminToken  *gocloak.JWT
	PublicKey   *rsa.PublicKey
	gocloak.GoCloak
	mx sync.RWMutex
}

func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	cli := gocloak.NewClient(cfg.Keycloak.BaseURL)

	var pk *rsa.PublicKey

	i, err := cli.GetIssuer(ctx, cfg.Keycloak.Realm)
	if err != nil {
		return nil, err
	}

	pbKeyBytes, err := base64.StdEncoding.DecodeString(*i.PublicKey)
	if err != nil {
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(pbKeyBytes)
	if err != nil {
		return nil, err
	}

	switch publicKey := pub.(type) {
	case *rsa.PublicKey:
		pk = publicKey
	default:
		return nil, fmt.Errorf("issuer return not rsa key")
	}

	return &Client{
		GoCloak:   cli,
		Cfg:       cfg,
		PublicKey: pk,
	}, nil
}

func (c *Client) setAdminToken(token *gocloak.JWT) {
	c.mx.Lock()
	c.adminToken = token
	c.mx.Unlock()
}

func (c *Client) getAdminToken() (gocloak.JWT, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if c.adminToken != nil {
		return *c.adminToken, true
	}

	return gocloak.JWT{}, false
}

func (c *Client) RefreshTokenOrRelogin(ctx context.Context, refreshToken, clientID, clientSecret, realm string) (*gocloak.JWT, error) {
	newAdminToken, err := c.RefreshToken(ctx, refreshToken, clientID, clientSecret, realm)
	if err == nil {
		return newAdminToken, nil
	}

	username := c.Cfg.Keycloak.Admin.Username
	password := c.Cfg.Keycloak.Admin.Password
	newAdminToken, err = c.LoginAdmin(ctx, username, password, realm)
	if err != nil {
		return nil, err
	}

	return newAdminToken, nil
}

func (c *Client) GetOrExtractAccessToken(ctx context.Context) (string, error) {
	if c.Cfg.Keycloak.Admin.Username == "" {
		return "", fmt.Errorf("no admin creds for auth cli")
	}

	username := c.Cfg.Keycloak.Admin.Username
	password := c.Cfg.Keycloak.Admin.Password
	realm := c.Cfg.Keycloak.Realm

	adminToken, exist := c.getAdminToken()
	if !exist {
		newAdminToken, err := c.LoginAdmin(ctx, username, password, realm)
		if err != nil {
			return "", err
		}
		c.setAdminToken(newAdminToken)

		return newAdminToken.AccessToken, nil
	}

	session, err := c.SessionFromToken(adminToken.AccessToken)
	if err != nil {
		validationError, ok := err.(*jwt.ValidationError)
		if ok {
			if validationError.Errors == jwt.ValidationErrorExpired {
				newAdminToken, errAuth := c.RefreshTokenOrRelogin(ctx, adminToken.RefreshToken, "admin-cli", "", realm)
				if errAuth != nil {
					return "", errAuth
				}

				c.setAdminToken(newAdminToken)

				return newAdminToken.AccessToken, nil
			}
		}
		return "", err
	}

	tokenExpire := time.Unix(session.ExpiresAt, 0)
	if time.Now().After(tokenExpire.Add(time.Minute * time.Duration(-1*TokenMinuteCorrectionFactor))) {
		newAdminToken, err := c.RefreshTokenOrRelogin(ctx, adminToken.RefreshToken, "admin-cli", "", realm)
		if err != nil {
			return "", err
		}

		c.setAdminToken(newAdminToken)

		return newAdminToken.AccessToken, nil
	}

	return adminToken.AccessToken, nil
}

func (c *Client) SessionFromToken(token string) (*Session, error) {
	if token == "" {
		return nil, ErrNoTokenFound
	}

	token = strings.Replace(token, "Bearer ", "", 1)
	session := new(Session)

	if _, err := jwt.ParseWithClaims(token, session, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return c.PublicKey, nil
	}); err != nil {
		return nil, err
	}

	if err := session.Valid(); err != nil {
		return nil, err
	}

	return session, nil
}

func (c *Client) ExtractUserByUsername(ctx context.Context, username string) (*User, error) {
	params := gocloak.GetUsersParams{
		Username: &username,
	}

	token, err := c.GetOrExtractAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	realm := c.Cfg.Keycloak.Realm
	users, err := c.GoCloak.GetUsers(ctx, token, realm, params)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNoUserFound
	}

	for _, val := range users {
		if *val.Username == username {
			return NewUserFromKeycloakModel(val), nil
		}
	}
	return nil, ErrNoUserFound
}
