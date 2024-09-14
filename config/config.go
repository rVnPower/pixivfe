// Global (Server-Wide) Settings

package config

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

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

	Token          []string `env:"PIXIVFE_TOKEN,required"` // may be multiple tokens. delimiter is ','
	LoadBalancing  string   `env:"PIXIVFE_TOKEN_LOAD_BALANCING,overwrite"` // 'round-robin' or 'random' 
	InDevelopment  bool     `env:"PIXIVFE_DEV"`
	UserAgent      string   `env:"PIXIVFE_USERAGENT,overwrite"`
	AcceptLanguage string   `env:"PIXIVFE_ACCEPTLANGUAGE,overwrite"`
	RequestLimit   int      `env:"PIXIVFE_REQUESTLIMIT"` // if 0, request limit is disabled

	ProxyServer_staging string  `env:"PIXIVFE_IMAGEPROXY,overwrite"`
	ProxyServer         url.URL // proxy server URL, may or may not contain authority part of the URL

	ProxyCheckInterval time.Duration `env:"PIXIVFE_PROXY_CHECK_INTERVAL,overwrite"`

	tokenIndex uint32 // Used for round-robin token selection
}

func (s *ServerConfig) LoadConfig() error {
	s.Version = "v2.8"
	log.Printf("PixivFE %s\n", s.Version)

	s.StartingTime = time.Now().UTC().Format("2006-01-02 15:04")

	// set default values with env:"...,overwrite"
	s.UserAgent = "Mozilla/5.0 (Windows NT 10.0; rv:123.0) Gecko/20100101 Firefox/123.0"
	s.AcceptLanguage = "en-US,en;q=0.5"
	s.ProxyServer_staging = BuiltinProxyUrl
	s.ProxyCheckInterval = 8 * time.Hour
	s.LoadBalancing = "round-robin"

	// load config from from env vars
	if err := envconfig.Process(context.Background(), s); err != nil {
		return err
	}

	if s.Port == "" && s.UnixSocket == "" {
		log.Fatalln("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
		return errors.New("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
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
	log.Printf("Load balancing method: %s\n", s.LoadBalancing)

	return nil
}

func (s *ServerConfig) GetToken() string {
	switch s.LoadBalancing {
	case "random":
		return s.getRandomToken()
	default:
		return s.getRoundRobinToken()
	}
}

func (s *ServerConfig) getRandomToken() string {
	return s.Token[rand.Intn(len(s.Token))]
}

func (s *ServerConfig) getRoundRobinToken() string {
	index := atomic.AddUint32(&s.tokenIndex, 1) % uint32(len(s.Token))
	return s.Token[index]
}
