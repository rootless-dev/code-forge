package bootstrap

import (
	"context"
	"github.com/carlosealves2/code-forge/oidc-service/configs"
	"github.com/carlosealves2/code-forge/oidc-service/internal/app"
	"github.com/carlosealves2/code-forge/oidc-service/internal/infra/cache"
	"github.com/carlosealves2/code-forge/oidc-service/internal/infra/http/handler"
	"github.com/carlosealves2/code-forge/oidc-service/internal/infra/http/server"
	"github.com/carlosealves2/code-forge/oidc-service/internal/usecase/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	_ "github.com/joho/godotenv/autoload"
	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"time"
)

func New(version, shortCommit string) *app.App {
	cfg := configs.NewConfigBuilder().WithEnv(version, shortCommit).Validate().Build()

	oidcConfig, provider := initOidc(cfg)
	cacheService := initCacheService(cfg)

	redirectUseCase := initRedirectUseCase(oidcConfig, cacheService)

	callbackUseCase := initCallbackUseCase(cfg, oidcConfig, provider, cacheService)

	httpServer := server.NewHttpServer(cfg.AppName(), version, shortCommit)

	setupMiddlewares(httpServer, cfg)

	server.AddRoute(httpServer, initHttpHandler(redirectUseCase, callbackUseCase))

	return app.New(cfg, httpServer)
}

func initOidc(cfg *configs.AppConfig) (*oauth2.Config, *oidc.Provider) {
	provider, err := oidc.NewProvider(context.Background(), cfg.OIDC().ProviderURL())
	if err != nil {
		log.Fatal().Err(err).Msg("error creating oidc provider")
	}

	return &oauth2.Config{
		ClientID:     cfg.OIDC().ClientID(),
		ClientSecret: cfg.OIDC().Secret(),
		RedirectURL:  cfg.OIDC().RedirectUrl(),

		Endpoint: provider.Endpoint(),

		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}, provider
}

func initCacheService(cfg *configs.AppConfig) cache.Cache {
	opt, err := redis.ParseURL(cfg.Redis().URI())
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing redis uri")
	}
	rdb := redis.NewClient(opt)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("error connecting to redis")
	}

	cacheClient, err := cache.New(cache.RedisCacheType, rdb)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating cache client")
	}
	return cacheClient
}

func initRedirectUseCase(oauth2Config *oauth2.Config, cacheClient cache.Cache) auth.RedirectUseCase {
	return auth.NewRedirectUseCase(&auth.RedirectUseCaseOptions{
		OAuth2Config: oauth2Config,
		CacheClient:  cacheClient,
	})
}

func initCallbackUseCase(
	cfg *configs.AppConfig,
	oauth2Config *oauth2.Config,
	provider *oidc.Provider,
	cacheClient cache.Cache) auth.CallbackOauth2UseCase {
	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.OIDC().ClientID()})

	return auth.NewCallbackOauth2UseCase(verifier, oauth2Config, cacheClient)
}

func initHttpHandler(redirectUseCase auth.RedirectUseCase, callbackUseCase auth.CallbackOauth2UseCase) *handler.OIDCHandler {
	return handler.NewOIDCHandler(&handler.OIDCHandlerDependencies{
		RedirectUseCase: redirectUseCase,
		CallbackUseCase: callbackUseCase,
	})
}

// setupMiddlewares configures middleware for the provided Fiber app, including logging, CORS, healthcheck, rate limiting, and profiling.
func setupMiddlewares(app *fiber.App, cfg *configs.AppConfig) {
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.HttpServer().AllowedOrigins(),
		AllowMethods: "GET",
	}))

	app.Use(healthcheck.New())

	app.Use(limiter.New(limiter.Config{
		Max: 5,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		Expiration: 30 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
		},
		Storage: nil,
	}))

	app.Use(pprof.New())
}
