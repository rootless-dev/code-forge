package server

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"strings"
)

func IsAuth(verifier *oidc.IDTokenVerifier) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := verifier.Verify(c.Context(), tokenStr)
		if err != nil {
			log.Error().Err(err).Msg("error verifying token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
				"hint":  err.Error(),
			})
		}
		return c.Next()
	}
}
