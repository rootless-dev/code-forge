package server

import "github.com/gofiber/fiber/v2"

func NewHttpServer() *fiber.App {
	return fiber.New()
}
