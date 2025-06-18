package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/carlosealves2/code-forge/oidc-service/internal/dto"
	"github.com/carlosealves2/code-forge/oidc-service/internal/infra/cache"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/phuslu/log"
	"golang.org/x/oauth2"
	"time"
)

// CallbackOauth2UseCase defines a contract for handling OAuth2 callback logic with input validation and token processing.
type CallbackOauth2UseCase interface {
	Execute(ctx context.Context, input *CallbackOauth2UseCaseInput) (*dto.Token, error)
}

// ImplCallbackOauth2UseCase is an implementation of the CallbackOauth2UseCase interface for handling OAuth2 callbacks.
// It uses an OIDC verifier, an OAuth2 configuration, and a cache client for state validation and token handling.
type ImplCallbackOauth2UseCase struct {
	verifier     OIDCVerifierInterface
	oauth2Config Oauth2ConfigInterface
	cacheClient  cache.Cache
}

// ValidationRule represents a function that performs a validation and returns an error if the validation fails.
type ValidationRule func() error

var (

	// ErrCallbackStateNotFound indicates that the specified state could not be found in the cache.
	ErrCallbackStateNotFound = errors.New("callback: state not found")

	// ErrCallbackInvalidState represents an error indicating that the provided state is invalid or cannot be unmarshalled properly.
	ErrCallbackInvalidState = errors.New("callback: invalid state")

	// ErrCallbackUserAgentMismatch indicates that the user agent in the callback request does not match the original state.
	ErrCallbackUserAgentMismatch = errors.New("callback: user agent does not match")

	// ErrCallbackIPMismatch indicates a mismatch between the expected and actual IP in a callback request validation.
	ErrCallbackIPMismatch = errors.New("callback: ip does not match")

	// ErrCallbackRequestExpired indicates that the callback request has expired based on a timestamp validation.
	ErrCallbackRequestExpired = errors.New("callback: request expired")

	// ErrCallbackExchangeCode occurs when there is an error exchanging an authorization code for an OAuth2 token.
	ErrCallbackExchangeCode = errors.New("callback: error exchanging code for token")

	// ErrCallbackVerifyIDToken represents an error that occurs during the verification of an ID token in the callback process.
	ErrCallbackVerifyIDToken = errors.New("callback: error verifying id token")

	// ErrCallbackRemoveState indicates an error occurred while attempting to remove a state from the cache.
	ErrCallbackRemoveState = errors.New("callback: error removing state from cache")

	// ErrCallbackIDTokenNotFound represents an error indicating that the ID token was not found during the callback process.
	ErrCallbackIDTokenNotFound = errors.New("callback: id token not found")
)

// NewCallbackOauth2UseCase creates a new instance of CallbackOauth2UseCase with OIDC verifier, OAuth2 configuration, and cache client.
func NewCallbackOauth2UseCase(verifier OIDCVerifierInterface, oauth2Config Oauth2ConfigInterface, cacheClient cache.Cache) CallbackOauth2UseCase {
	return &ImplCallbackOauth2UseCase{
		verifier:     verifier,
		oauth2Config: oauth2Config,
		cacheClient:  cacheClient,
	}
}

// CallbackOauth2UseCaseInput represents the input required to handle an OAuth2 callback request.
// It includes fields like authorization code, state identifier, client IP, user agent, and request timestamp.
type CallbackOauth2UseCaseInput struct {
	Code      string
	State     string
	IP        string
	UserAgent string
	Timestamp time.Time
}

// CallbackOauth2UseCaseOutput represents the output of an OAuth2 callback use case containing access and refresh tokens.
type CallbackOauth2UseCaseOutput struct {
	AccessToken  string
	RefreshToken string
}

// Execute handles the OAuth2 callback process, performs validations, and returns the access and refresh tokens.
func (u *ImplCallbackOauth2UseCase) Execute(ctx context.Context, input *CallbackOauth2UseCaseInput) (*dto.Token, error) {
	cachedState, err := u.retrieveStateFromCache(ctx, input.State)
	if err != nil {
		return nil, err
	}

	if err = u.validateRequest(ctx, cachedState, input); err != nil {
		if !errors.Is(err, ErrCallbackRemoveState) {
			return nil, err
		}
	}

	oauth2Token, err := u.exchangeToken(ctx, input.Code)
	if err != nil {
		return nil, err
	}

	rawIdToken, err := u.getRawIDTokenFromOAuth2Token(oauth2Token)
	if err != nil {
		return nil, err
	}

	_, err = u.validateIDToken(ctx, rawIdToken, cachedState.Nonce)
	if err != nil {
		return nil, err
	}

	return &dto.Token{
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
	}, nil
}

// retrieveStateFromCache retrieves the state data from the cache using the provided state identifier.
// Returns the parsed state data or an error if retrieval or unmarshalling fails.
func (u *ImplCallbackOauth2UseCase) retrieveStateFromCache(ctx context.Context, state string) (*dto.StateData, error) {
	rawState, err := u.cacheClient.Get(ctx, state)
	if err != nil {
		log.Error().Err(err).Msg("error retrieving state from cache")
		return nil, ErrCallbackStateNotFound
	}

	var stateData dto.StateData

	err = json.Unmarshal(rawState, &stateData)
	if err != nil {
		log.Error().Err(err).Msg("error unmarshalling state data")
		return nil, ErrCallbackInvalidState
	}

	return &stateData, nil
}

// validateRequest validates the OAuth2 callback request using cached state and input; removes state from cache upon success.
func (u *ImplCallbackOauth2UseCase) validateRequest(ctx context.Context, cachedState *dto.StateData, input *CallbackOauth2UseCaseInput) error {
	rules := []ValidationRule{
		func() error {
			if cachedState.UserAgent != input.UserAgent {
				return ErrCallbackUserAgentMismatch
			}
			return nil
		},

		func() error {
			if cachedState.IP != input.IP {
				return ErrCallbackIPMismatch
			}
			return nil
		},

		func() error {
			if input.Timestamp.Sub(cachedState.Timestamp) > 5*time.Minute {
				return ErrCallbackRequestExpired
			}
			return nil
		},
	}

	for _, rule := range rules {
		if err := rule(); err != nil {
			return err
		}
	}

	err := u.cacheClient.Remove(ctx, cachedState.ID)
	if err != nil {
		log.Error().Err(err).Msg("error removing state from cache")
		return ErrCallbackRemoveState
	}

	return nil
}

// exchangeToken exchanges an authorization code for an OAuth2 token using the provided context and code parameters.
func (u *ImplCallbackOauth2UseCase) exchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	oauth2Token, err := u.oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Error().Err(err).Msg("error exchanging code for token")
		return nil, ErrCallbackExchangeCode
	}

	return oauth2Token, err
}

// validateIDToken verifies the provided ID token, checks its nonce against the expected value, and returns the parsed ID token.
func (u *ImplCallbackOauth2UseCase) validateIDToken(ctx context.Context, token string, expectedNonce string) (*oidc.IDToken, error) {
	idToken, err := u.verifier.Verify(ctx, token)
	if err != nil {
		log.Error().Err(err).Msg("error verifying id token")
		return nil, ErrCallbackVerifyIDToken
	}

	if idToken.Nonce != expectedNonce {
		log.Error().Msg("nonce does not match")
		return nil, ErrCallbackInvalidState
	}

	return idToken, nil
}

// getRawIDTokenFromOAuth2Token extracts the raw ID token as a string from the OAuth2 token's extra data field.
func (u *ImplCallbackOauth2UseCase) getRawIDTokenFromOAuth2Token(token *oauth2.Token) (string, error) {
	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", ErrCallbackIDTokenNotFound
	}
	return rawIdToken, nil
}
