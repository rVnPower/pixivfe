package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/proxy_checker"
	"codeberg.org/vnpower/pixivfe/v2/server/audit"
	"codeberg.org/vnpower/pixivfe/v2/server/handlers"
	"codeberg.org/vnpower/pixivfe/v2/server/template"
)

func main() {
	config.GlobalConfig.LoadConfig()
	audit.Init(config.GlobalConfig.InDevelopment)
	template.Init(config.GlobalConfig.InDevelopment)

	// Initialize and start the proxy checker
	ctx_timeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	proxy_checker.InitializeProxyChecker(ctx_timeout)
	handlers.InitializeRateLimiter()

	router := handlers.DefineRoutes()
	// the first middleware is the most outer / first executed one
	router.Use(handlers.ProvideUserContext) // needed for everything else
	router.Use(handlers.LogRequest)         // all pages need this
	router.Use(handlers.SetPrivacyHeaders)  // all pages need this
	router.Use(handlers.HandleError)        // if the inner handler fails, this shows the error page instead
	router.Use(handlers.RateLimitRequest)

	// watch and compile sass when in development mode
	if config.GlobalConfig.InDevelopment {
		go run_sass()
	}

	// Listen
	err := http.Serve(chooseListener(), router)
	if err != http.ErrServerClosed {
		log.Print(err)
	}
}

func run_sass() {
	cmd := exec.Command("sass", "--watch", "assets/css")
	cmd.Stdout = os.Stderr // Sass quirk
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pdeathsig: syscall.SIGHUP}
	runtime.LockOSThread() // Go quirk https://github.com/golang/go/issues/27505
	err := cmd.Run()
	if err != nil {
		log.Print(fmt.Errorf("when running sass: %w", err))
	}
}

func chooseListener() net.Listener {
	var l net.Listener
	if config.GlobalConfig.UnixSocket != "" {
		ln, err := net.Listen("unix", config.GlobalConfig.UnixSocket)
		if err != nil {
			panic(err)
		}
		l = ln
		log.Printf("Listening on domain socket %v\n", config.GlobalConfig.UnixSocket)
	} else {
		addr := config.GlobalConfig.Host + ":" + config.GlobalConfig.Port
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			log.Panicf("failed to listen: %v", err)
		}
		l = ln
		addr = ln.Addr().String()
		log.Printf("Listening on http://%v/\n", addr)
	}
	return l
}
