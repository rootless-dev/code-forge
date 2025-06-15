package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/carlosealves2/code-forge/oidc-service/internal/core/entity"
	"github.com/carlosealves2/code-forge/oidc-service/internal/infra/cache"
	"github.com/google/uuid"
	"github.com/phuslu/log"
	"golang.org/x/oauth2"
	"time"
)

type RedirectUseCase interface {
	Execute(context.Context, *RedirectUseCaseInput) string
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

type RedirectUseCaseInput struct {
	UserAgent string
	IP        string
	Timestamp time.Time
}

func (u *ImplRedirectUseCase) Execute(ctx context.Context, input *RedirectUseCaseInput) string {
	stateData := u.generateState(input)

	rawStateData, err := json.Marshal(stateData)
	if err != nil {
		log.Error().Err(err).Msg("error serializing state")
		return ""
	}

	err = u.saveCache(ctx, stateData.ID, rawStateData)
	if err != nil {
		return ""
	}

	return u.oauth2Config.AuthCodeURL(stateData.ID)
}

func (u *ImplRedirectUseCase) generateState(input *RedirectUseCaseInput) *entity.StateData {
	return &entity.StateData{
		ID:        uuid.New().String(),
		Timestamp: input.Timestamp,
		IP:        input.IP,
		UserAgent: input.UserAgent,
	}
}

func (u *ImplRedirectUseCase) saveCache(ctx context.Context, key string, data []byte) error {
	err := u.cacheClient.Set(ctx, key, data, 5*time.Minute)
	if err != nil {
		log.Error().Err(err).Msg("error storing state in cache")
		return ErrStoreStateInCache
	}

	return nil
}
