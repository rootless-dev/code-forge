package app

import (
	"fmt"
	"github.com/carlosealves2/code-forge/oidc-service/configs"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
)

type App struct {
	cfg        *configs.AppConfig
	httpServer *fiber.App
}

func New(cfg *configs.AppConfig, httpServer *fiber.App) *App {
	return &App{
		httpServer: httpServer,
		cfg:        cfg,
	}
}

func (a *App) Run() {
	addrs := fmt.Sprintf("%s:%s", a.cfg.HttpServer().Host(), a.cfg.HttpServer().Port())
	if err := a.httpServer.Listen(addrs); err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}
