package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/carlosealves2/code-forge/oidc-service/internal/core/entity"
	"github.com/carlosealves2/code-forge/oidc-service/internal/infra/cache"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

func setupRedisContainerAndGetServiceURL(t *testing.T) (string, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	testcontainers.CleanupContainer(t, redisC)

	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return "", err
	}

	return redisC.Endpoint(ctx, "")
}

func setupFakeFaceOIDCConfig() *oauth2.Config {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("../../../testassets/well-known.json")
		if err != nil {
			log.Error().Err(err).Msg("error opening file")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			log.Error().Err(err).Msg("error reading file")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		content = bytes.ReplaceAll(content, []byte("localhost"), []byte(r.Host))

		var contentMap map[string]interface{}

		json.Unmarshal(content, &contentMap)

		json.NewEncoder(w).Encode(contentMap)
	}))
	provider, err := oidc.NewProvider(context.Background(), testServer.URL+"/realms/code-forge")
	if err != nil {
		log.Fatal().Err(err)
		return nil
	}

	return &oauth2.Config{
		ClientID:     "fake-face-client-id",
		ClientSecret: "fake-face-client-secret",
		RedirectURL:  "http://localhost:8080/callback",

		Endpoint: provider.Endpoint(),

		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
}

func TestNewRedirectUseCase(t *testing.T) {
	redisUri, err := setupRedisContainerAndGetServiceURL(t)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	rdb := redis.NewClient(&redis.Options{Addr: redisUri})

	cacheClient, err := cache.New(cache.RedisCacheType, rdb)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	openidConfig := setupFakeFaceOIDCConfig()

	options := &RedirectUseCaseOptions{
		OAuth2Config: openidConfig,
		CacheClient:  cacheClient,
	}

	redirectUseCase := NewRedirectUseCase(options)

	input := &RedirectUseCaseInput{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:120.0) Gecko/20100101 Firefox/120.0",
		IP:        "127.0.0.1",
		Timestamp: time.Now().UTC(),
	}

	redirectUri := redirectUseCase.Execute(context.Background(), input)
	if redirectUri == "" {
		t.Errorf("expected a valid redirect uri, got %s", redirectUri)
	}

	url, err := url.Parse(redirectUri)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	state := url.Query().Get("state")
	if state == "" {
		t.Errorf("expected a valid state, got %s", state)
		return
	}

	data, err := cacheClient.Get(context.Background(), state)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	var stateData entity.StateData

	err = json.Unmarshal(data, &stateData)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	if stateData.UserAgent != input.UserAgent {
		t.Errorf("expected user agent to be the same as the value stored in cache, got %s", stateData.UserAgent)
		return
	}

	if stateData.IP != input.IP {
		t.Errorf("expected ip to be the same as the value stored in cache, got %s", stateData.IP)
		return
	}

	if stateData.Timestamp != input.Timestamp {
		t.Errorf("expected timestamp to be the same as the value stored in cache, got %s", stateData.Timestamp)
	}
}

type MockCacheClient struct {
	GetFunc func(ctx context.Context, key string) ([]byte, error)
	SetFunc func(ctx context.Context, key string, value []byte, expirationTime time.Duration) error
}

func (m MockCacheClient) Get(ctx context.Context, key string) ([]byte, error) {
	return m.GetFunc(ctx, key)
}

func (m MockCacheClient) Set(ctx context.Context, key string, value []byte, expirationTime time.Duration) error {
	return m.SetFunc(ctx, key, value, expirationTime)
}

func TestSaveCacheShouldReturnError(t *testing.T) {
	m := MockCacheClient{
		SetFunc: func(ctx context.Context, key string, value []byte, expirationTime time.Duration) error {
			return errors.New("error")
		},
	}

	redirectUseCase := &ImplRedirectUseCase{
		oauth2Config: setupFakeFaceOIDCConfig(),
		cacheClient:  m,
	}

	err := redirectUseCase.saveCache(context.Background(), "test", []byte("test"))
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

func TestImplRedirectUseCase_Execute_ShouldReturnAnEmptyString(t *testing.T) {
	m := MockCacheClient{
		SetFunc: func(ctx context.Context, key string, value []byte, expirationTime time.Duration) error {
			return errors.New("error")
		},
	}

	redirectUseCase := &ImplRedirectUseCase{
		oauth2Config: &oauth2.Config{},
		cacheClient:  m,
	}

	input := &RedirectUseCaseInput{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:120.0) Gecko/20100101 Firefox/120.0",
		IP:        "127.0.0.1",
		Timestamp: time.Now().UTC(),
	}

	redirectUri := redirectUseCase.Execute(context.Background(), input)
	if redirectUri != "" {
		t.Errorf("expected an empty string, got %s", redirectUri)
	}
}
