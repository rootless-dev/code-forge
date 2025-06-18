package auth

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"testing"
)

type MockTokenSource struct {
	TokenFunc func() (*oauth2.Token, error)
}

func (m *MockTokenSource) Token() (*oauth2.Token, error) {
	return m.TokenFunc()
}

func TestNewRefreshTokenUseCase(t *testing.T) {
	instance := NewRefreshTokenUseCase(new(MockOAuth2Config))
	if instance == nil {
		t.Errorf("expected valid instance, got nil")
		return
	}
}

func TestRefreshTokenUseCase_Success(t *testing.T) {
	configMock := &MockOAuth2Config{
		TokenSourceFunc: func(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
			return &MockTokenSource{
				TokenFunc: func() (*oauth2.Token, error) {
					return &oauth2.Token{
						AccessToken:  "access-token-abc",
						RefreshToken: "refresh-token-xyz",
					}, nil
				},
			}
		},
	}

	instance := NewRefreshTokenUseCase(configMock)
	if instance == nil {
		t.Errorf("expected valid instance, got nil")
		return
	}

	ctx := context.Background()
	input := &RefreshTokenUseCaseInput{
		RefreshToken: "old-refresh-token",
	}

	out, err := instance.Execute(ctx, input)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	if out == nil {
		t.Errorf("expected valid output, got nil")
		return
	}

	if out.AccessToken == "" {
		t.Errorf("expected valid access token, got %s", out.AccessToken)
		return
	}

	if out.RefreshToken == "" {
		t.Errorf("expected valid refresh token, got %s", out.RefreshToken)
		return
	}
}

func TestRefreshTokenUseCase_NewTokenError(t *testing.T) {
	configMock := &MockOAuth2Config{
		TokenSourceFunc: func(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
			return &MockTokenSource{
				TokenFunc: func() (*oauth2.Token, error) {
					return nil, errors.New("error")
				},
			}
		},
	}

	instance := NewRefreshTokenUseCase(configMock)

	ctx := context.Background()
	input := &RefreshTokenUseCaseInput{
		RefreshToken: "old-refresh-token",
	}

	out, err := instance.Execute(ctx, input)
	if err == nil {
		t.Errorf("expected an error, got %v", err)
	}

	if out != nil {
		t.Errorf("expected nil output, got %v", out)
	}

}
