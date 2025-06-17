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
)

func New(version, shortCommit string) *app.App {
	cfg := configs.NewConfigBuilder().WithEnv(version, shortCommit).Validate().Build()

	oidcConfig, provider := initOidc(cfg)
	cacheService := initCacheService(cfg)

	initRedirectUseCase(oidcConfig, cacheService)

	httpServer := server.NewHttpServer(cfg.AppName(), version, shortCommit)

	setupMiddlewares(httpServer, cfg)

	server.AddRoute(httpServer, initHttpHandler(initRedirectUseCase(oidcConfig, cacheService)))

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

func initHttpHandler(redirectUseCase auth.RedirectUseCase) *handler.OIDCHandler {
	return handler.NewOIDCHandler(&handler.OIDCHandlerDependencies{
		RedirectUseCase: redirectUseCase,
	})
}

func setupMiddlewares(app *fiber.App, cfg *configs.AppConfig) {
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.HttpServer().AllowedOrigins(),
		AllowMethods: "GET",
	}))

	app.Use(healthcheck.New())

	app.Use(limiter.New())

	app.Use(pprof.New())
}
