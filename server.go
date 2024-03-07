package xgo

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	auth "github.com/anoaland/xgo/auth"
)

type AuthManager interface {
	GetCurrentUser(ctx *fiber.Ctx) *auth.AppUser
}

type WebServer struct {
	App  *fiber.App
	Auth *auth.WebAuthManager
}

type XRouter struct {
	fiber.Router
	ws WebServer
}

func (xr XRouter) WithAuth(prefix string) *XRouter {
	return &XRouter{
		xr.ws.WithAuth(xr, prefix),
		xr.ws,
	}
}

func (xr XRouter) XGroup(prefix string) *XRouter {
	return &XRouter{
		xr.Group(prefix),
		xr.ws,
	}
}

func New() *WebServer {
	app := fiber.New()

	return &WebServer{
		App: app,
	}
}

func (s *WebServer) UseAuth(client auth.WebAuthClient, bearerTokenConfig *auth.BearerTokenMiddlewareConfig) {
	s.Auth = auth.NewWebAuthManager(client, bearerTokenConfig)
}

func (s *WebServer) XGroup(prefix string) *XRouter {
	return &XRouter{
		s.App.Group(prefix),
		*s,
	}
}

func (s *WebServer) WithAuth(r fiber.Router, group string) fiber.Router {
	return r.Group(group, s.Auth.AuthGuardMiddleware)
}

func (server *WebServer) Run(port int, onShutdown func() error) {

	// Listen for syscall signals for process to interrupt/quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		_ = <-c
		fmt.Println("\r\nGracefully shutting down...")
		_ = server.App.Shutdown()
	}()
	// see: https://adrianhesketh.com/2021/05/28/templ-hot-reload-with-air/
	addr := fmt.Sprintf(":%d", port)
	err := server.App.Listen(addr)

	if err != nil {
		log.Fatal(err)
	}

	err = onShutdown()
	if err != nil {
		log.Fatal(err)
	}
}
