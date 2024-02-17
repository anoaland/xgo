package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const USER_LOCAL_KEY = "x-user"

type WebAuthClient interface {
	GetUserFromToken(token string) (*AppUser, error)
}

type WebAuthManager struct {
	bearerTokenConfig *BearerTokenMiddlewareConfig
	client            WebAuthClient
}

func NewWebAuthManager(client WebAuthClient, opts *BearerTokenMiddlewareConfig) *WebAuthManager {
	config := &BearerTokenMiddlewareConfig{
		BodyKey:    "access_token",
		HeaderKey:  "Bearer",
		QueryKey:   "access_token",
		RequestKey: "token",
	}

	if opts != nil {
		if len(opts.BodyKey) > 0 {
			config.BodyKey = opts.BodyKey
		}

		if len(opts.HeaderKey) > 0 {
			config.HeaderKey = opts.HeaderKey
		}

		if len(opts.QueryKey) > 0 {
			config.QueryKey = opts.QueryKey
		}

		if len(opts.RequestKey) > 0 {
			config.RequestKey = opts.RequestKey
		}
	}

	return &WebAuthManager{bearerTokenConfig: config, client: client}
}

func (m *WebAuthManager) AuthGuardMiddleware(ctx *fiber.Ctx) error {
	var token string

	// get bearer token from request authorization header
	headerValue := ctx.Get("authorization")

	if len(headerValue) > 0 {
		components := strings.SplitN(headerValue, " ", 2)

		if len(components) == 2 && components[0] == m.bearerTokenConfig.HeaderKey {
			token = components[1]
		}
	} else {
		// get bearer token from query parameter
		queryValue := ctx.Query(m.bearerTokenConfig.QueryKey)

		if len(queryValue) > 0 {
			token = queryValue
		}
	}
	//
	// else, we might want to get token from Body or Request Parameters
	//

	user, err := m.client.GetUserFromToken(token)

	if err != nil {
		return err
	}

	ctx.Locals(USER_LOCAL_KEY, user)
	if user == nil {
		return ctx.SendStatus(401)
	}

	return ctx.Next()
}

func (m *WebAuthManager) User(ctx *fiber.Ctx) *AppUser {

	appUser, ok := ctx.Locals(USER_LOCAL_KEY).(*AppUser)

	if !ok {
		return nil
	}

	return appUser
}
