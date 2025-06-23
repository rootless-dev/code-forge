package auth

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Oauth2ConfigInterface defines a contract for exchanging an authorization code for an OAuth2 token using specified options.
type Oauth2ConfigInterface interface {
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource
}

// OIDCVerifierInterface defines the method to validate and parse an OpenID Connect ID token.
type OIDCVerifierInterface interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

type OIDCProviderInterface interface {
	UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*oidc.UserInfo, error)
}
