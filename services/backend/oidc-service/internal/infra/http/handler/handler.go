package handler

import (
	"github.com/carlosealves2/code-forge/oidc-service/internal/usecase/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"net/http"
	"strings"
	"time"
)

// OIDCHandler provides methods to handle OpenID Connect (OIDC) redirection and callback processes for authentication flows.
type OIDCHandler struct {
	redirectUseCase     auth.RedirectUseCase
	callbackUseCase     auth.CallbackOauth2UseCase
	refreshTokenUseCase auth.RefreshTokenUseCase
}

// OIDCHandlerDependencies contains the dependencies required for handling OIDC operations, including Redirect and Callback use cases.
type OIDCHandlerDependencies struct {
	RedirectUseCase     auth.RedirectUseCase
	CallbackUseCase     auth.CallbackOauth2UseCase
	RefreshTokenUseCase auth.RefreshTokenUseCase
}

// NewOIDCHandler creates and returns a new instance of OIDCHandler with the provided dependencies.
func NewOIDCHandler(deps *OIDCHandlerDependencies) *OIDCHandler {
	return &OIDCHandler{
		redirectUseCase:     deps.RedirectUseCase,
		callbackUseCase:     deps.CallbackUseCase,
		refreshTokenUseCase: deps.RefreshTokenUseCase,
	}
}

// Redirect handles the initiation of the OAuth2 login flow by redirecting the user to the appropriate OAuth2 provider URL.
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

// Callback handles the callback from the OAuth2 provider, processes authorization code and state, and returns a JSON response.
func (o *OIDCHandler) Callback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			log.Error().Msg("code is empty")
			return c.Status(fiber.StatusBadRequest).SendString("code is empty")
		}

		state := c.Query("state")
		if state == "" {
			log.Error().Msg("state is empty")
			return c.Status(fiber.StatusBadRequest).SendString("state is empty")
		}

		out, err := o.callbackUseCase.Execute(c.Context(), &auth.CallbackOauth2UseCaseInput{
			Code:      code,
			State:     state,
			IP:        c.IP(),
			UserAgent: c.Get("User-Agent"),
			Timestamp: time.Now().UTC(),
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{
					"error": err.Error(),
					"code":  http.StatusInternalServerError,
				})
		}

		return c.JSON(out)
	}
}

func (o *OIDCHandler) RefreshToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		refreshToken := c.Get("Authorization")
		if refreshToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
				"code":  http.StatusUnauthorized,
			})
		}

		refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
		if len(refreshToken) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
				"code":  http.StatusUnauthorized,
			})
		}

		out, err := o.refreshTokenUseCase.Execute(c.Context(), &auth.RefreshTokenUseCaseInput{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
				"code":  http.StatusInternalServerError,
			})
		}

		return c.JSON(out)
	}
}
