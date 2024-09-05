// Global (Server-Wide) Settings

package config

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/url"
	"strings"
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

	Token          []string `env:"PIXIVFE_TOKEN,required"` // may be multiple tokens. delimiter is ','    maybe add some testing?
	InDevelopment  bool     `env:"PIXIVFE_DEV"`
	UserAgent      string   `env:"PIXIVFE_USERAGENT,overwrite"`
	AcceptLanguage string   `env:"PIXIVFE_ACCEPTLANGUAGE,overwrite"`
	RequestLimit   int      `env:"PIXIVFE_REQUESTLIMIT"` // if 0, request limit is disabled

	ProxyServer_staging string  `env:"PIXIVFE_IMAGEPROXY,overwrite"`
	ProxyServer         url.URL // proxy server, may contain prefix as well

	ProxyCheckInterval_staging int           `env:"PIXIVFE_PROXY_CHECK_INTERVAL,overwrite"`
	ProxyCheckInterval         time.Duration // Proxy check interval
}

func (s *ServerConfig) LoadConfig() error {
	s.Version = "v2.8"
	log.Printf("PixivFE %s\n", s.Version)

	s.StartingTime = time.Now().UTC().Format("2006-01-02 15:04")

	// set default values with env:"...,overwrite"
	s.UserAgent = "Mozilla/5.0 (Windows NT 10.0; rv:123.0) Gecko/20100101 Firefox/123.0"
	s.AcceptLanguage = "en-US,en;q=0.5"
	s.ProxyServer_staging = BuiltinProxyUrl
	s.ProxyCheckInterval_staging = 480

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
	log.Printf("Set %s to: %s\n", "proxy server", &s.ProxyServer)

	s.ProxyCheckInterval = time.Duration(s.ProxyCheckInterval_staging) * time.Minute
	log.Printf("Proxy check interval set to: %v\n", s.ProxyCheckInterval)

	return nil
}

func GetRandomDefaultToken() string {
	defaultToken := GlobalConfig.Token[rand.Intn(len(GlobalConfig.Token))]

	return defaultToken
}
