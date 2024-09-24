// Package config provides global server-wide settings.
package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/server/token_manager"
	"github.com/sethvargo/go-envconfig"
)

var GlobalConfig ServerConfig

// REVISION stores the current version's revision information
var REVISION string = ""

const (
	unknownRevision = "unknown"
	revisionFormat  = "date-hash[+dirty]"
)

type ServerConfig struct {
	Version      string
	Revision     string
	RevisionDate string
	RevisionHash string
	IsDirty      bool
	StartingTime string // used in /about page

	Host string `env:"PIXIVFE_HOST"`

	// One of the two is required
	Port       string `env:"PIXIVFE_PORT"`
	UnixSocket string `env:"PIXIVFE_UNIXSOCKET"`

	RepoURL string `env:"PIXIVFE_REPO_URL,overwrite"` // used in /about page

	Token               []string `env:"PIXIVFE_TOKEN,required"` // may be multiple tokens. delimiter is ','
	TokenManager        *token_manager.TokenManager
	TokenLoadBalancing  string        `env:"PIXIVFE_TOKEN_LOAD_BALANCING,overwrite"`
	TokenMaxRetries     int           `env:"PIXIVFE_TOKEN_MAX_RETRIES,overwrite"`
	TokenBaseTimeout    time.Duration `env:"PIXIVFE_TOKEN_BASE_TIMEOUT,overwrite"`
	TokenMaxBackoffTime time.Duration `env:"PIXIVFE_TOKEN_MAX_BACKOFF_TIME,overwrite"`

	// API request level backoff settings
	APIMaxRetries     int           `env:"PIXIVFE_API_MAX_RETRIES,overwrite"`
	APIBaseTimeout    time.Duration `env:"PIXIVFE_API_BASE_TIMEOUT,overwrite"`
	APIMaxBackoffTime time.Duration `env:"PIXIVFE_API_MAX_BACKOFF_TIME,overwrite"`

	AcceptLanguage string `env:"PIXIVFE_ACCEPTLANGUAGE,overwrite"`
	RequestLimit   int    `env:"PIXIVFE_REQUESTLIMIT"` // if 0, request limit is disabled

	ProxyServer_staging string  `env:"PIXIVFE_IMAGEPROXY,overwrite"`
	ProxyServer         url.URL // proxy server URL, may or may not contain authority part of the URL

	ProxyCheckEnabled  bool          `env:"PIXIVFE_PROXY_CHECK_ENABLED,overwrite"`
	ProxyCheckInterval time.Duration `env:"PIXIVFE_PROXY_CHECK_INTERVAL,overwrite"`

	// Development options
	InDevelopment        bool   `env:"PIXIVFE_DEV"`
	ResponseSaveLocation string `env:"PIXIVFE_RESPONSE_SAVE_LOCATION,overwrite"`
}

// parseRevision extracts date, hash, and dirty status from the revision string
func parseRevision(revision string) (date, hash string, isDirty bool) {
	if revision == "" {
		return unknownRevision, unknownRevision, false
	}

	isDirty = strings.HasSuffix(revision, "+dirty")
	if isDirty {
		revision = strings.TrimSuffix(revision, "+dirty")
	}

	parts := strings.Split(revision, "-")
	if len(parts) == 2 {
		return parts[0], parts[1], isDirty
	}
	return unknownRevision, revision, isDirty
}

// validateURL checks if the given URL is valid
func validateURL(urlString string, urlType string) (*url.URL, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	if (parsedURL.Scheme == "") != (parsedURL.Host == "") {
		return nil, fmt.Errorf("%s URL is invalid: %s. Please specify e.g. https://example.com", urlType, urlString)
	}
	if strings.HasSuffix(parsedURL.Path, "/") {
		return nil, fmt.Errorf("%s URL path (%s) cannot end in /: %s. PixivFE does not support this now", urlType, parsedURL.Path, urlString)
	}
	return parsedURL, nil
}

func (s *ServerConfig) LoadConfig() error {
	s.Version = "v2.9"
	s.Revision = REVISION

	s.RevisionDate, s.RevisionHash, s.IsDirty = parseRevision(REVISION)

	if REVISION == "" {
		log.Printf("[WARNING] REVISION is not set. Continuing with unknown revision information.\n")
	} else if s.RevisionDate == unknownRevision {
		log.Printf("[WARNING] REVISION format is invalid: %s. Expected format '%s'. Continuing with full revision as hash.\n", REVISION, revisionFormat)
	}

	log.Printf("PixivFE %s, revision %s\n", s.Version, s.Revision)

	s.StartingTime = time.Now().UTC().Format("2006-01-02 15:04")

	// set default values with env:"...,overwrite"
	s.RepoURL = "https://codeberg.org/VnPower/PixivFE"

	s.AcceptLanguage = "en-US,en;q=0.5"
	s.ProxyServer_staging = BuiltinProxyUrl
	s.ProxyCheckEnabled = true
	s.ProxyCheckInterval = 8 * time.Hour
	s.TokenLoadBalancing = "round-robin"
	s.TokenMaxRetries = 5
	s.TokenBaseTimeout = 1000 * time.Millisecond
	s.TokenMaxBackoffTime = 32000 * time.Millisecond

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
		return errors.New("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set")
	}

	// a check for tokens
	if len(s.Token) < 1 {
		log.Fatalln("PIXIVFE_TOKEN has to be set. Visit https://pixivfe-docs.pages.dev/hosting/hosting-pixivfe for more details.")
		return errors.New("PIXIVFE_TOKEN has to be set. Visit https://pixivfe-docs.pages.dev/hosting/hosting-pixivfe for more details")
	}

	// Validate proxy server URL
	proxyURL, err := validateURL(s.ProxyServer_staging, "Proxy server")
	if err != nil {
		log.Printf("[WARNING] Invalid proxy server URL: %v. Falling back to built-in proxy URL.", err)
		proxyURL, _ = url.Parse(BuiltinProxyUrl) // We know this is valid
	} else {
		log.Printf("Proxy server set to: %s\n", proxyURL.String())
	}
	s.ProxyServer = *proxyURL
	log.Printf("Proxy server set to: %s\n", s.ProxyServer.String())
	log.Printf("Proxy check interval set to: %v\n", s.ProxyCheckInterval)

	// Validate repo URL
	repoURL, err := validateURL(s.RepoURL, "Repo")
	if err != nil {
		log.Printf("[WARNING] Invalid repo URL: %v. Using default repo URL.", err)
		s.RepoURL = "https://codeberg.org/VnPower/PixivFE" // Use a default value
	} else {
		s.RepoURL = repoURL.String()
		log.Printf("Repo URL set to: %s\n", s.RepoURL)
	}

	// Validate TokenLoadBalancing
	switch s.TokenLoadBalancing {
	case "round-robin", "random", "least-recently-used":
		// Valid options
	default:
		log.Printf("[WARNING] Invalid PIXIVFE_TOKEN_LOAD_BALANCING value: %s. Defaulting to 'round-robin'.\n", s.TokenLoadBalancing)
		s.TokenLoadBalancing = "round-robin"
	}

	// Initialize TokenManager
	s.TokenManager = token_manager.NewTokenManager(s.Token, s.TokenMaxRetries, s.TokenBaseTimeout, s.TokenMaxBackoffTime, s.TokenLoadBalancing)
	log.Printf("Token manager initialized with %d tokens\n", len(s.Token))
	log.Printf("Token manager settings: Max retries: %d, Base timeout: %v, Max backoff time: %v\n", s.TokenMaxRetries, s.TokenBaseTimeout, s.TokenMaxBackoffTime)
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
