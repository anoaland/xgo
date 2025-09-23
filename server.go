package xgo

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anoaland/xgo/internal"
	"github.com/gofiber/fiber/v2"

	auth "github.com/anoaland/xgo/auth"
)

type AuthManager interface {
	GetCurrentUser(ctx *fiber.Ctx) interface{}
}

// type WebServerErrorHandler = func(err xgoErrors.XgoError)

type WebServer struct {
	App  *fiber.App
	Auth *auth.WebAuthManager
	// errorHandler *WebServerErrorHandler
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

func New(config ...fiber.Config) *WebServer {
	var appConfig fiber.Config
	if len(config) > 0 {
		appConfig = config[0]
	} else {
		appConfig = fiber.Config{}
	}

	if appConfig.ErrorHandler == nil {
		appConfig.ErrorHandler = DefaultErrorHandler()
	}

	app := fiber.New(appConfig)

	return &WebServer{
		App: app,
	}
}

func (s *WebServer) UseAuth(client auth.WebAuthClient, bearerTokenConfig *auth.BearerTokenMiddlewareConfig) {
	s.Auth = auth.NewWebAuthManager(client, bearerTokenConfig)
}

// func (s *WebServer) UseErrorHandler(fn WebServerErrorHandler) {
// 	s.errorHandler = &fn
// }

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

func (server *WebServer) RunOnAddress(addr string, onShutdown func() error) {

	// Listen for syscall signals for process to interrupt/quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		_ = <-c
		fmt.Println("\r\nGracefully shutting down...")
		_ = server.App.Shutdown()
	}()

	err := server.App.Listen(addr)

	if err != nil {
		log.Fatal(err)
	}

	err = onShutdown()
	if err != nil {
		log.Fatal(err)
	}
}

func (server *WebServer) LoggerContext(ctx *fiber.Ctx) context.Context {
	// Get the per-request logger with request_id context
	requestLogger := GetRequestLogger(ctx)
	if requestLogger == nil {
		// Fallback to basic context if no request logger is available
		return context.WithValue(context.Background(), internal.FiberContextKey, ctx)
	}

	// Create a context with the request logger embedded
	loggerCtx := requestLogger.WithContext(context.Background())

	// Also store the fiber context for the GORM logger to access
	loggerCtx = context.WithValue(loggerCtx, internal.FiberContextKey, ctx)

	return loggerCtx
}
