// Global (Server-Wide) Settings

package config

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/server/token_manager"
	"github.com/sethvargo/go-envconfig"
)

var GlobalConfig ServerConfig

type ServerConfig struct {
	Version      string
	StartingTime string // used in /about page

	Host string `env:"PIXIVFE_HOST"`

	// One of the two is required
	Port       string `env:"PIXIVFE_PORT"`
	UnixSocket string `env:"PIXIVFE_UNIXSOCKET"`

	Token              []string `env:"PIXIVFE_TOKEN,required"` // may be multiple tokens. delimiter is ','
	TokenManager       *token_manager.TokenManager
	TokenLoadBalancing string        `env:"PIXIVFE_TOKEN_LOAD_BALANCING,overwrite"`
	MaxRetries         int           `env:"PIXIVFE_MAX_RETRIES,overwrite"`
	BaseTimeout        time.Duration `env:"PIXIVFE_BASE_TIMEOUT,overwrite"`
	MaxBackoffTime     time.Duration `env:"PIXIVFE_MAX_BACKOFF_TIME,overwrite"`

	// API request level backoff settings
	APIMaxRetries     int           `env:"PIXIVFE_API_MAX_RETRIES,overwrite"`
	APIBaseTimeout    time.Duration `env:"PIXIVFE_API_BASE_TIMEOUT,overwrite"`
	APIMaxBackoffTime time.Duration `env:"PIXIVFE_API_MAX_BACKOFF_TIME,overwrite"`

	UserAgent      string `env:"PIXIVFE_USERAGENT,overwrite"`
	AcceptLanguage string `env:"PIXIVFE_ACCEPTLANGUAGE,overwrite"`
	RequestLimit   int    `env:"PIXIVFE_REQUESTLIMIT"` // if 0, request limit is disabled

	ProxyServer_staging string  `env:"PIXIVFE_IMAGEPROXY,overwrite"`
	ProxyServer         url.URL // proxy server URL, may or may not contain authority part of the URL

	ProxyCheckInterval time.Duration `env:"PIXIVFE_PROXY_CHECK_INTERVAL,overwrite"`

	// Development options
	InDevelopment  bool   `env:"PIXIVFE_DEV"`
	ResponseSaveLocation string `env:"PIXIVFE_RESPONSE_SAVE_LOCATION,overwrite"`
}

func (s *ServerConfig) LoadConfig() error {
	s.Version = "v2.9"
	log.Printf("PixivFE %s\n", s.Version)

	s.StartingTime = time.Now().UTC().Format("2006-01-02 15:04")

	// set default values with env:"...,overwrite"
	s.UserAgent = "Mozilla/5.0 (Windows NT 10.0; rv:123.0) Gecko/20100101 Firefox/123.0"
	s.AcceptLanguage = "en-US,en;q=0.5"
	s.ProxyServer_staging = BuiltinProxyUrl
	s.ProxyCheckInterval = 8 * time.Hour
	s.TokenLoadBalancing = "round-robin"
	s.MaxRetries = 5
	s.BaseTimeout = 1000 * time.Millisecond
	s.MaxBackoffTime = 32000 * time.Millisecond

	s.APIMaxRetries = 3
	s.APIBaseTimeout = 500 * time.Millisecond
	s.APIMaxBackoffTime = 8000 * time.Millisecond

	s.ResponseSaveLocation = "/tmp/pixivfe/responses"

	// load config from from env vars
	if err := envconfig.Process(context.Background(), s); err != nil {
		return err
	}

	if s.Port == "" && s.UnixSocket == "" {
		log.Fatalln("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
		return errors.New("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
	}

	// a check for tokens
	if len(s.Token) < 1 {
		log.Fatalln("PIXIVFE_TOKEN has to be set. Visit https://pixivfe-docs.pages.dev/hosting/hosting-pixivfe for more details.")
		return errors.New("PIXIVFE_TOKEN has to be set. Visit https://pixivfe-docs.pages.dev/hosting/hosting-pixivfe for more details.")
	}

	{ // validate proxy server
		proxyUrl, err := url.Parse(s.ProxyServer_staging)
		if err != nil {
			panic(err)
		}
		s.ProxyServer = *proxyUrl
		if (proxyUrl.Scheme == "") != (proxyUrl.Host == "") {
			log.Panicf("proxy server url is weird: %s\nPlease specify e.g. https://example.com", proxyUrl.String())
		}
		if strings.HasSuffix(proxyUrl.Path, "/") {
			log.Panicf("proxy server path (%s) has cannot end in /: %s\nPixivFE does not support this now, sorry", proxyUrl.Path, proxyUrl.String())
		}
	}
	log.Printf("Proxy server set to: %s\n", s.ProxyServer.String())
	log.Printf("Proxy check interval set to: %v\n", s.ProxyCheckInterval)

	// Validate TokenLoadBalancing
	switch s.TokenLoadBalancing {
	case "round-robin", "random", "least-recently-used":
		// Valid options
	default:
		log.Printf("Invalid PIXIVFE_TOKEN_LOAD_BALANCING value: %s. Defaulting to 'round-robin'.\n", s.TokenLoadBalancing)
		s.TokenLoadBalancing = "round-robin"
	}

	// Initialize TokenManager
	s.TokenManager = token_manager.NewTokenManager(s.Token, s.MaxRetries, s.BaseTimeout, s.MaxBackoffTime, s.TokenLoadBalancing)
	log.Printf("Token manager initialized with %d tokens\n", len(s.Token))
	log.Printf("Token manager settings: Max retries: %d, Base timeout: %v, Max backoff time: %v\n", s.MaxRetries, s.BaseTimeout, s.MaxBackoffTime)
	log.Printf("Token load balancing method: %s\n", s.TokenLoadBalancing)

	log.Printf("API request backoff settings: Max retries: %d, Base timeout: %v, Max backoff time: %v\n", s.APIMaxRetries, s.APIBaseTimeout, s.APIMaxBackoffTime)

	// Only print ResponseSaveLocation if InDevelopment is set
	if s.InDevelopment {
		log.Printf("Response save location: %s\n", s.ResponseSaveLocation)
	}

	return nil
}

func (s *ServerConfig) GetToken() string {
	token := s.TokenManager.GetToken()
	if token == nil {
		log.Println("[WARNING] All tokens are timed out. Using the first available token.")
		return s.Token[0]
	}
	return token.Value
}
