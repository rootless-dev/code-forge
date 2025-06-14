package configs

import (
	_ "github.com/carlosealves2/code-forge/oidc-service/internal/bootstrap/logging"
	"github.com/phuslu/log"
	"net/url"
	"os"
)

type ConfigBuilder struct {
	cfg *AppConfig
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (b *ConfigBuilder) WithEnv(version, shortCommit string) *ConfigBuilder {
	b.cfg = &AppConfig{
		appName:     getEnv("APP_NAME", "oidc-service"),
		version:     version,
		shortCommit: shortCommit,
		env:         getEnv("APP_ENV", "development"),
		httpServerConfig: &HttpServerConfig{
			host:           getEnv("HTTP_SERVER_HOST", "localhost"),
			port:           getEnv("HTTP_SERVER_PORT", "3000"),
			allowedOrigins: getEnv("HTTP_SERVER_ALLOWED_ORIGINS", "*"),
		},
		oidcConfig: &OIDCConfig{
			providerURL: getEnv("OIDC_PROVIDER_URL", ""),
			clientID:    getEnv("OIDC_CLIENT_ID", ""),
			secret:      getEnv("OIDC_SECRET", ""),
			redirectUrl: getEnv("OIDC_REDIRECT_URL", "http://localhost:3000/oauth2/v1/callback"),
		},
		redisConfig: &RedisConfig{
			uri: getEnv("REDIS_URI", "redis://localhost:6379"),
		},
	}

	return b
}

func (b *ConfigBuilder) Validate() *ConfigBuilder {
	if b.cfg == nil {
		log.Fatal().Msg("configuration has not been initialized")
	}

	if b.cfg.oidcConfig.providerURL == "" {
		log.Fatal().Msg("OIDC_PROVIDER_URL is required")
	}

	if b.cfg.oidcConfig.clientID == "" {
		log.Fatal().Msg("OIDC_CLIENT_ID is required")
	}

	if b.cfg.oidcConfig.secret == "" {
		log.Fatal().Msg("OIDC_SECRET is required")
	}

	if b.cfg.env == "production" || b.cfg.env == "homolog" {
		callback, err := url.Parse(b.cfg.OIDC().redirectUrl)
		if err != nil {
			log.Fatal().Err(err).Msg("invalid redirect url")
		}

		if callback.Scheme != "https" {
			log.Fatal().Msg("redirect url must be https in production and homolog environments")
		}

		redisUriParsed, err := url.Parse(b.cfg.redisConfig.URI())
		if err != nil {
			log.Fatal().Err(err).Msg("invalid redis uri")
		}

		userInfo := redisUriParsed.User
		if userInfo == nil {
			log.Fatal().Msg("redis uri must have user info in production and homolog environments")
		}

		username := userInfo.Username()
		password, hasPassword := userInfo.Password()

		if username == "" || !hasPassword || password == "" {
			log.Fatal().Msg("redis uri must have user info in production and homolog environments")
		}

		if b.cfg.httpServerConfig.AllowedOrigins() == "*" {
			log.Fatal().Msg("http server allowed origins must be a valid url in production and homolog environments")
		}
	}

	return b
}

func (b *ConfigBuilder) Build() *AppConfig {
	return b.cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
