package auth

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/anoaland/xgo/auth"
)

type KeycloakAppUser struct {
	*gocloak.UserInfo
}

func (u *KeycloakAppUser) AsAppUser() *auth.AppUser {
	return &auth.AppUser{
		Username: *u.PreferredUsername,
	}
}
