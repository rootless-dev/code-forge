package auth

import (
	"context"
	"errors"
	"github.com/carlosealves2/code-forge/oidc-service/internal/infra/cache"
	"github.com/google/uuid"
	"github.com/phuslu/log"
	"golang.org/x/oauth2"
	"time"
)

type RedirectUseCase interface {
	Execute(context.Context) string
}
type ImplRedirectUseCase struct {
	oauth2Config *oauth2.Config
	cacheClient  cache.Cache
}

type RedirectUseCaseOptions struct {
	OAuth2Config *oauth2.Config
	CacheClient  cache.Cache
}

var (
	ErrStoreStateInCache = errors.New("error storing state in cache")
)

func NewRedirectUseCase(options *RedirectUseCaseOptions) RedirectUseCase {
	return &ImplRedirectUseCase{
		oauth2Config: options.OAuth2Config,
		cacheClient:  options.CacheClient,
	}
}

func (u *ImplRedirectUseCase) Execute(ctx context.Context) string {
	state := u.generateState()

	err := u.saveCache(ctx, state, []byte(state))
	if err != nil {
		return ""
	}

	return u.oauth2Config.AuthCodeURL(state)
}

func (u *ImplRedirectUseCase) generateState() string {
	return uuid.NewString()
}

func (u *ImplRedirectUseCase) saveCache(ctx context.Context, key string, data []byte) error {
	err := u.cacheClient.Set(ctx, key, data, 5*time.Minute)
	if err != nil {
		log.Error().Err(err).Msg("error storing state in cache")
		return ErrStoreStateInCache
	}

	return nil
}
