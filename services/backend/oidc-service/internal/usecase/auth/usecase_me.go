package auth

import (
	"context"
	"github.com/phuslu/log"
	"golang.org/x/oauth2"
)

type MeUseCase interface {
	GetUserInfo(ctx context.Context, input *MeUseCaseInput) (map[string]any, error)
}

type ImplMeUseCase struct {
	verifier     OIDCVerifierInterface
	provider     OIDCProviderInterface
	oauth2Config Oauth2ConfigInterface
}

func NewMeUseCase(provider OIDCProviderInterface, verifier OIDCVerifierInterface, oauth2Config Oauth2ConfigInterface) MeUseCase {
	return &ImplMeUseCase{
		verifier:     verifier,
		provider:     provider,
		oauth2Config: oauth2Config,
	}
}

type MeUseCaseInput struct {
	AccessToken string
}

func (u *ImplMeUseCase) GetUserInfo(ctx context.Context, input *MeUseCaseInput) (map[string]any, error) {
	tokenSource := u.getTokenSource(ctx, input.AccessToken)

	userInfo, err := u.provider.UserInfo(ctx, tokenSource)
	if err != nil {
		log.Error().Err(err).Msg("error getting user info")
		return nil, err
	}

	var claims map[string]any

	if err := userInfo.Claims(&claims); err != nil {
		log.Error().Err(err).Msg("error parsing claims")
		return nil, err
	}

	return claims, nil
}

func (u *ImplMeUseCase) getTokenSource(ctx context.Context, accessToken string) oauth2.TokenSource {
	return u.oauth2Config.TokenSource(ctx, &oauth2.Token{
		AccessToken: accessToken,
	})
}
