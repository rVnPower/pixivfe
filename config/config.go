// Global (Server-Wide) Settings

package config

import (
	"errors"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var GlobalServerConfig ServerConfig

type ServerConfig struct {
	// Required
	Token []string

	ProxyServer url.URL // proxy server, may contain prefix as well

	// can be left empty
	Host string

	// One of two is required
	Port       string
	UnixSocket string

	UserAgent      string
	AcceptLanguage string
	RequestLimit   int // if 0, request limit is disabled

	StartingTime  string
	Version       string
	InDevelopment bool

	ProxyCheckInterval time.Duration // Proxy check interval
}

func (s *ServerConfig) LoadConfig() error {
	s.setVersion()

	CollectAllEnv()

	token, hasToken := LookupEnv("PIXIVFE_TOKEN")
	if !hasToken {
		log.Fatalln("PIXIVFE_TOKEN is required, but was not set.")
		return errors.New("PIXIVFE_TOKEN is required, but was not set.\n")
	}
	// TODO Maybe add some testing?
	s.Token = strings.Split(token, ",")

	port, hasPort := LookupEnv("PIXIVFE_PORT")
	socket, hasSocket := LookupEnv("PIXIVFE_UNIXSOCKET")
	if !hasPort && !hasSocket {
		log.Fatalln("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
		return errors.New("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
	}
	s.Port = port
	s.UnixSocket = socket

	_, hasDev := LookupEnv("PIXIVFE_DEV")
	s.InDevelopment = hasDev

	s.Host = GetEnv("PIXIVFE_HOST")

	s.UserAgent = GetEnv("PIXIVFE_USERAGENT")

	s.AcceptLanguage = GetEnv("PIXIVFE_ACCEPTLANGUAGE")

	s.SetRequestLimit(GetEnv("PIXIVFE_REQUESTLIMIT"))

	s.SetProxyServer(GetEnv("PIXIVFE_IMAGEPROXY"))

	s.SetProxyCheckInterval(GetEnv("PIXIVFE_PROXY_CHECK_INTERVAL"))

	AnnounceAllEnv()

	s.setStartingTime()

	return nil
}

func (s *ServerConfig) SetProxyServer(v string) {
	proxyUrl, err := url.Parse(v)
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

func (s *ServerConfig) SetRequestLimit(v string) {
	if v == "" {
		s.RequestLimit = 0
	} else {
		t, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		s.RequestLimit = t
	}
}

func (s *ServerConfig) SetProxyCheckInterval(v string) {
	const defaultInterval = 480
	if v == "" {
		s.ProxyCheckInterval = defaultInterval * time.Minute
	} else {
		minutes, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("Invalid PIXIVFE_PROXY_CHECK_INTERVAL value: %s. Using default of %d minutes.\n", v, defaultInterval)
			s.ProxyCheckInterval = defaultInterval * time.Minute
		} else {
			s.ProxyCheckInterval = time.Duration(minutes) * time.Minute
		}
	}
	log.Printf("Proxy check interval set to: %v\n", s.ProxyCheckInterval)
}

func (s *ServerConfig) setStartingTime() {
	s.StartingTime = time.Now().UTC().Format("2006-01-02 15:04")
	log.Printf("Set starting time to: %s\n", s.StartingTime)
}

func (s *ServerConfig) setVersion() {
	s.Version = "v2.7.1"
	log.Printf("PixivFE %s\n", s.Version)
}

func GetRandomDefaultToken() string {
	defaultToken := GlobalServerConfig.Token[rand.Intn(len(GlobalServerConfig.Token))]

	return defaultToken
}
