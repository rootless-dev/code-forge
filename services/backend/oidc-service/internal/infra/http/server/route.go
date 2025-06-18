package server

import (
	"github.com/carlosealves2/code-forge/oidc-service/internal/infra/http/handler"
	"github.com/gofiber/fiber/v2"
)

func AddRoute(app *fiber.App, oidcHandler *handler.OIDCHandler) {
	app.Get("/sso/login", oidcHandler.Redirect())
	app.Get("/sso/callback", oidcHandler.Callback())
}
