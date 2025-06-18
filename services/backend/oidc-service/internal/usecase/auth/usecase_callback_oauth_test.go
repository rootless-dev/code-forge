package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/carlosealves2/code-forge/oidc-service/internal/dto"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"testing"
	"time"
)

type MockVerifier struct {
	VerifyFunc func(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

func (m *MockVerifier) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	return m.VerifyFunc(ctx, rawIDToken)
}

type MockOAuth2Config struct {
	ExchangeFunc    func(ctx context.Context, code string) (*oauth2.Token, error)
	TokenSourceFunc func(ctx context.Context, t *oauth2.Token) oauth2.TokenSource
}

func (m *MockOAuth2Config) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return m.ExchangeFunc(ctx, code)
}

func (m *MockOAuth2Config) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	return m.TokenSourceFunc(ctx, t)
}

func TestCallbackUseCase_Success(t *testing.T) {
	state := "state-xyz"
	nonce := "nonce-abc123"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			data, _ := json.Marshal(dto.StateData{
				ID:        state,
				Nonce:     nonce,
				IP:        "127.0.0.1",
				UserAgent: "test-agent",
				Timestamp: time.Now(),
			})
			return data, nil
		},
		RemoveFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}

	verifier := &MockVerifier{
		VerifyFunc: func(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
			return &oidc.IDToken{Nonce: nonce}, nil
		},
	}
	config := &MockOAuth2Config{
		ExchangeFunc: func(ctx context.Context, code string) (*oauth2.Token, error) {
			token := &oauth2.Token{
				AccessToken:  "access-token-abc",
				RefreshToken: "valid-refresh-token",
			}
			return token.WithExtra(map[string]any{
				"id_token": "fake-id-token",
			}), nil
		},
	}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "test-agent",
		Timestamp: time.Now(),
	}

	output, err := useCase.Execute(context.Background(), input)
	if err != nil {
		t.Error("expected no error, got ", err)
		return
	}

	if output.AccessToken != "access-token-abc" {
		t.Errorf("expected access token, got %s", output.AccessToken)
	}
}

func TestCallbackUseCase_NonceMismatch(t *testing.T) {
	state := "state-xyz"
	correctNonce := "nonce-correct"
	invalidNonce := "nonce-invalid"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			data, _ := json.Marshal(dto.StateData{
				ID:        state,
				Nonce:     correctNonce,
				IP:        "127.0.0.1",
				UserAgent: "test-agent",
				Timestamp: time.Now(),
			})
			return data, nil
		},
		RemoveFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}

	verifier := &MockVerifier{
		VerifyFunc: func(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
			return &oidc.IDToken{Nonce: invalidNonce}, nil
		},
	}
	config := &MockOAuth2Config{
		ExchangeFunc: func(ctx context.Context, code string) (*oauth2.Token, error) {
			token := &oauth2.Token{
				AccessToken:  "access-token-abc",
				RefreshToken: "valid-refresh-token",
			}
			return token.WithExtra(map[string]any{
				"id_token": "fake-id-token",
			}), nil
		},
	}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "test-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected nonce mismatch error, got nil")
	}

	if !errors.Is(err, ErrCallbackInvalidState) {
		t.Errorf("expected ErrCallbackInvalidState, got %s", err)
	}
}

func TestCallbackUseCase_StateExpired(t *testing.T) {
	state := "state-xyz"
	nonce := "nonce-abc123"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			data, _ := json.Marshal(dto.StateData{
				ID:        state,
				Nonce:     nonce,
				IP:        "127.0.0.1",
				UserAgent: "test-agent",
				Timestamp: time.Now().Add(-10 * time.Minute),
			})
			return data, nil
		},
		RemoveFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}

	verifier := &MockVerifier{}
	config := &MockOAuth2Config{}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "test-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected expiration error, got nil")
	}
}

func TestCallbackUseCase_UserAgentMismatch(t *testing.T) {
	state := "state-xyz"
	nonce := "nonce-abc123"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			data, _ := json.Marshal(dto.StateData{
				ID:        state,
				Nonce:     nonce,
				IP:        "127.0.0.1",
				UserAgent: "expected-agent",
				Timestamp: time.Now(),
			})
			return data, nil
		},
		RemoveFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}

	verifier := &MockVerifier{}
	config := &MockOAuth2Config{}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "wrong-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected user agent mismatch error, got nil")
	}
}

func TestNewCallbackOauth2UseCase(t *testing.T) {
	instance := NewCallbackOauth2UseCase(new(MockVerifier), new(MockOAuth2Config), new(MockCacheClient))

	if instance == nil {
		t.Errorf("expected valid instance, got nil")
		return
	}
}

func TestCallbackUseCase_IpMismatch(t *testing.T) {
	state := "state-xyz"
	nonce := "nonce-abc123"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			data, _ := json.Marshal(dto.StateData{
				ID:        state,
				Nonce:     nonce,
				IP:        "127.0.0.1",
				UserAgent: "expected-agent",
				Timestamp: time.Now(),
			})
			return data, nil
		},
		RemoveFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}

	verifier := &MockVerifier{}
	config := &MockOAuth2Config{}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "192.168.1.1",
		UserAgent: "expected-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected ip mismatch error, got nil")
	}
}

func TestCallbackUseCase_RemoveStateError(t *testing.T) {
	state := "state-xyz"
	nonce := "nonce-abc123"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			data, _ := json.Marshal(dto.StateData{
				ID:        state,
				Nonce:     nonce,
				IP:        "127.0.0.1",
				UserAgent: "expected-agent",
				Timestamp: time.Now(),
			})
			return data, nil
		},
		RemoveFunc: func(ctx context.Context, key string) error {
			return errors.New("error")
		},
	}

	verifier := &MockVerifier{
		VerifyFunc: func(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
			return &oidc.IDToken{Nonce: nonce}, nil
		},
	}
	config := &MockOAuth2Config{
		ExchangeFunc: func(ctx context.Context, code string) (*oauth2.Token, error) {
			token := &oauth2.Token{
				AccessToken:  "access-token-abc",
				RefreshToken: "valid-refresh-token",
			}
			return token.WithExtra(map[string]any{
				"id_token": "fake-id-token",
			}), nil
		},
	}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "expected-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err != nil {
		t.Error("expected a no error, but got", err)
	}
}

func TestCallbackUseCase_RetrieveStateError(t *testing.T) {
	state := "state-xyz"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			return nil, errors.New("error")
		},
	}

	verifier := &MockVerifier{}
	config := &MockOAuth2Config{}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "expected-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error but got nil")
	}

	if !errors.Is(err, ErrCallbackStateNotFound) {
		t.Errorf("expected ErrCallbackStateNotFound, got %s", err)
	}
}

func TestCallbackUseCase_ParseStateError(t *testing.T) {
	state := "state-xyz"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			return []byte(`invalid-state`), nil
		},
	}

	verifier := &MockVerifier{}
	config := &MockOAuth2Config{}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "expected-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error but got nil")
	}

	if !errors.Is(err, ErrCallbackInvalidState) {
		t.Errorf("expected ErrCallbackInvalidState, got %s", err)
	}
}

func TestCallbackUseCase_ExchangeError(t *testing.T) {
	state := "state-xyz"
	nonce := "nonce-abc123"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			data, _ := json.Marshal(dto.StateData{
				ID:        state,
				Nonce:     nonce,
				IP:        "127.0.0.1",
				UserAgent: "expected-agent",
				Timestamp: time.Now(),
			})
			return data, nil
		},
		RemoveFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}

	verifier := &MockVerifier{}
	config := &MockOAuth2Config{
		ExchangeFunc: func(ctx context.Context, code string) (*oauth2.Token, error) {
			return nil, errors.New("error")
		},
	}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "expected-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error but got nil")
	}

	if !errors.Is(err, ErrCallbackExchangeCode) {
		t.Errorf("expected ErrCallbackExchangeCode, got %s", err)
	}
}

func TestCallbackUseCase_VerifierError(t *testing.T) {
	state := "state-xyz"
	nonce := "nonce-abc123"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			data, _ := json.Marshal(dto.StateData{
				ID:        state,
				Nonce:     nonce,
				IP:        "127.0.0.1",
				UserAgent: "expected-agent",
				Timestamp: time.Now(),
			})
			return data, nil
		},
		RemoveFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}

	verifier := &MockVerifier{
		VerifyFunc: func(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
			return nil, errors.New("error")
		},
	}
	config := &MockOAuth2Config{
		ExchangeFunc: func(ctx context.Context, code string) (*oauth2.Token, error) {
			token := &oauth2.Token{
				AccessToken:  "access-token-abc",
				RefreshToken: "refresh-token-xyz",
			}
			return token.WithExtra(map[string]any{
				"id_token": "fake-id-token",
			}), nil
		},
	}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "expected-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error but got nil")
	}

	if !errors.Is(err, ErrCallbackVerifyIDToken) {
		t.Errorf("expected ErrCallbackVerifyIDToken, got %s", err)
	}
}

func TestCallbackUseCase_ExtraRawIdTokenError(t *testing.T) {
	state := "state-xyz"
	nonce := "nonce-abc123"

	mockCache := &MockCacheClient{
		GetFunc: func(ctx context.Context, key string) ([]byte, error) {
			data, _ := json.Marshal(dto.StateData{
				ID:        state,
				Nonce:     nonce,
				IP:        "127.0.0.1",
				UserAgent: "expected-agent",
				Timestamp: time.Now(),
			})
			return data, nil
		},
		RemoveFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}

	verifier := &MockVerifier{
		VerifyFunc: func(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
			return nil, errors.New("error")
		},
	}
	config := &MockOAuth2Config{
		ExchangeFunc: func(ctx context.Context, code string) (*oauth2.Token, error) {
			token := &oauth2.Token{
				AccessToken:  "access-token-abc",
				RefreshToken: "refresh-token-xyz",
			}
			return token, nil
		},
	}

	useCase := &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: config,
		cacheClient:  mockCache,
	}

	input := &CallbackOauth2UseCaseInput{
		Code:      "dummy-code",
		State:     state,
		IP:        "127.0.0.1",
		UserAgent: "expected-agent",
		Timestamp: time.Now(),
	}

	_, err := useCase.Execute(context.Background(), input)
	if err == nil {
		t.Error("expected error but got nil")
	}

	if !errors.Is(err, ErrCallbackIDTokenNotFound) {
		t.Errorf("expected ErrCallbackIDTokenNotFound, got %s", err)
	}
}
