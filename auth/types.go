package auth

type AppUser struct {
	Username string
}

// BearerTokenMiddlewareConfig holds the configuration of the middleware. It is completely optional
// and should only be provided if your application uses token keys that are not
// RFC6750-compliant.
type BearerTokenMiddlewareConfig struct {
	// BodyKey defines the key to use when searching for the bearer token inside the
	// request's body.
	// Optional. Default: "access_token".
	BodyKey string

	// HeaderKey defines the prefix of the Authorization header's value, used when
	// searching for the bearer token inside the request's headers.
	// Optional. Default: "Bearer".
	HeaderKey string

	// QueryKey defines the key to use when searching for the bearer token inside the
	// request's query parameters.
	// Optional. Default: "access_token".
	QueryKey string

	// RequestKey defines the name of the local variable that will be created in the
	// request's context, which will contain the bearer token extracted from the
	// request.
	// Optional. Default: "token".
	RequestKey string
}
