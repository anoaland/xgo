package auth

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/anoaland/xgo/auth"
)

type KeycloakWebAuthClient struct {
	kk           *gocloak.GoCloak
	url          string
	realm        string
	clientId     string
	clientSecret string
}

func New(url string, realm string, clientId string, clientSecret string) *KeycloakWebAuthClient {
	kk := gocloak.NewClient(url)

	return &KeycloakWebAuthClient{
		kk:           kk,
		url:          url,
		realm:        realm,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (c KeycloakWebAuthClient) GetUserFromToken(token string) (*auth.AppUser, error) {
	user, err := c.kk.GetUserInfo(context.Background(), token, c.realm)

	if err != nil {
		return nil, err
	}

	var kuser = KeycloakAppUser{user}

	return kuser.AsAppUser(), nil
}

func (c KeycloakWebAuthClient) Login(ctx context.Context, usernameOrEmail string, password string) (*gocloak.JWT, error) {
	return c.kk.Login(ctx, c.clientId, c.clientSecret, c.realm, usernameOrEmail, password)
}

func (c KeycloakWebAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	return c.kk.RefreshToken(ctx, refreshToken, c.clientId, c.clientSecret, c.realm)
}

func (c KeycloakWebAuthClient) Register(ctx context.Context, user gocloak.User, password string) (*string, error) {

	token, err := c.kk.GetToken(ctx, c.realm, gocloak.TokenOptions{
		ClientID:     &c.clientId,
		ClientSecret: &c.clientSecret,
		GrantType:    gocloak.StringP("client_credentials"),
	})

	// this works as well
	// token, err := auth.client.LoginAdmin(c, "dd-admin", "dd-admin", realm)

	if err != nil {
		return nil, err
	}

	userToken := token.AccessToken

	userId, err := c.kk.CreateUser(ctx, userToken, c.realm, user)

	if err != nil {
		return nil, err
	}

	err = c.kk.SetPassword(ctx, userToken, userId, c.realm, password, false)
	if err != nil {
		return nil, err
	}

	return &userId, nil
}
