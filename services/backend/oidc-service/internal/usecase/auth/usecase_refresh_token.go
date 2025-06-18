package auth

import (
	"context"
	"errors"
	"github.com/carlosealves2/code-forge/oidc-service/internal/dto"
	"github.com/phuslu/log"
	"golang.org/x/oauth2"
)

// RefreshTokenUseCase defines a contract for refreshing OAuth2 tokens using a refresh token and returns new token details.
type RefreshTokenUseCase interface {
	Execute(ctx context.Context, input *RefreshTokenUseCaseInput) (*dto.Token, error)
}

// ErrRefreshTokenGenerateNewToken is returned when there is an error in generating a new refresh token.
var (
	ErrRefreshTokenGenerateNewToken = errors.New("refresh token: error generating new token")
)

// ImplRefreshTokenUseCase is an implementation of RefreshTokenUseCase for handling token refresh logic.
// It uses Oauth2ConfigInterface to perform token exchange and generate new tokens based on refresh tokens.
type ImplRefreshTokenUseCase struct {
	oauth2Config Oauth2ConfigInterface
}

// NewRefreshTokenUseCase creates a new instance of RefreshTokenUseCase using the provided Oauth2ConfigInterface implementation.
func NewRefreshTokenUseCase(oauth2Config Oauth2ConfigInterface) RefreshTokenUseCase {
	return &ImplRefreshTokenUseCase{
		oauth2Config: oauth2Config,
	}
}

// RefreshTokenUseCaseInput represents the input required for refreshing an OAuth2 token.
type RefreshTokenUseCaseInput struct {
	RefreshToken string
}

// Execute refreshes an OAuth2 token using the provided refresh token and returns a new access and refresh token if successful.
func (u *ImplRefreshTokenUseCase) Execute(ctx context.Context, input *RefreshTokenUseCaseInput) (*dto.Token, error) {
	tokenSource := u.tokenSource(ctx, input.RefreshToken)

	newToken, err := u.newToken(tokenSource)
	if err != nil {
		return nil, err
	}

	return &dto.Token{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	}, nil
}

// tokenSource retrieves an OAuth2 token source for refreshing tokens using the given context and refresh token.
func (u *ImplRefreshTokenUseCase) tokenSource(ctx context.Context, refreshToken string) oauth2.TokenSource {
	return u.oauth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
}

// newToken retrieves a new OAuth2 token using the provided token source and returns it or an error if token generation fails.
func (u *ImplRefreshTokenUseCase) newToken(tokenSource oauth2.TokenSource) (*oauth2.Token, error) {
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while refreshing token")
		return nil, ErrRefreshTokenGenerateNewToken
	}

	return newToken, nil
}
