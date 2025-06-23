package server

import (
	"github.com/carlosealves2/code-forge/oidc-service/internal/infra/http/handler"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
)

func AddRoute(app *fiber.App, verifier *oidc.IDTokenVerifier, oidcHandler *handler.OIDCHandler) {
	sso := app.Group("/sso")
	sso.Get("/login", oidcHandler.Redirect())
	sso.Get("/callback", oidcHandler.Callback())
	sso.Get("/refresh", oidcHandler.RefreshToken())
	private := sso.Group("/private", IsAuth(verifier))
	private.Get("/me", oidcHandler.Me())
}
