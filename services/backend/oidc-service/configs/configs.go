package configs

type AppConfig struct {
	appName          string
	version          string
	shortCommit      string
	env              string
	host             string
	port             string
	oidcConfig       *OIDCConfig
	redisConfig      *RedisConfig
	httpServerConfig *HttpServerConfig
}

type OIDCConfig struct {
	providerURL string
	clientID    string
	secret      string
	redirectUrl string
}

type RedisConfig struct {
	uri string
}

type HttpServerConfig struct {
	host           string
	port           string
	allowedOrigins string
}

func (a *AppConfig) AppName() string               { return a.appName }
func (a *AppConfig) Version() string               { return a.version }
func (a *AppConfig) ShortCommit() string           { return a.shortCommit }
func (a *AppConfig) Env() string                   { return a.env }
func (a *AppConfig) Host() string                  { return a.host }
func (a *AppConfig) Port() string                  { return a.port }
func (a *AppConfig) OIDC() *OIDCConfig             { return a.oidcConfig }
func (a *AppConfig) Redis() *RedisConfig           { return a.redisConfig }
func (a *AppConfig) HttpServer() *HttpServerConfig { return a.httpServerConfig }

func (o *OIDCConfig) ProviderURL() string { return o.providerURL }
func (o *OIDCConfig) ClientID() string    { return o.clientID }
func (o *OIDCConfig) Secret() string      { return o.secret }
func (o *OIDCConfig) RedirectUrl() string { return o.redirectUrl }

func (r *RedisConfig) URI() string { return r.uri }

func (h *HttpServerConfig) Host() string           { return h.host }
func (h *HttpServerConfig) Port() string           { return h.port }
func (h *HttpServerConfig) AllowedOrigins() string { return h.allowedOrigins }
