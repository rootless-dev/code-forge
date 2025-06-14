package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func NewHttpServer(appName, version, shortCommit string) *fiber.App {
	fullAppName := fmt.Sprintf("%s %s (%s)", appName, version, shortCommit)
	return fiber.New(fiber.Config{
		AppName: fullAppName,
	})
}
