package handler

import (
	"github.com/carlosealves2/code-forge/oidc-service/internal/usecase/auth"
	"github.com/gofiber/fiber/v2"
	"time"
)

type OIDCHandler struct {
	redirectUseCase auth.RedirectUseCase
}

type OIDCHandlerDependencies struct {
	RedirectUseCase auth.RedirectUseCase
}

func NewOIDCHandler(deps *OIDCHandlerDependencies) *OIDCHandler {
	return &OIDCHandler{
		redirectUseCase: deps.RedirectUseCase,
	}
}

func (o *OIDCHandler) Redirect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		redirectUri := o.redirectUseCase.Execute(c.Context(), &auth.RedirectUseCaseInput{
			UserAgent: c.Get("User-Agent"),
			IP:        c.IP(),
			Timestamp: time.Now().UTC(),
		})
		if redirectUri == "" {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Redirect(redirectUri, fiber.StatusFound)
	}
}
