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

	"codeberg.org/vnpower/pixivfe/v2/audit"
	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/handlers"
	"codeberg.org/vnpower/pixivfe/v2/template"
)

func main() {
	config.GlobalConfig.LoadConfig()
	audit.Init(config.GlobalConfig.InDevelopment)
	template.Init(config.GlobalConfig.InDevelopment)

	// Initialize and start the proxy checker
	ctx_timeout, cancel := context.WithTimeout(context.Background(), config.ProxyCheckerTimeout)
	defer cancel()
	config.InitializeProxyChecker(ctx_timeout)
	handlers.InitializeRateLimiter()

	router := handlers.DefineRoutes()

	main_handler := func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
		handlers.ErrorHandler(w, r)
	}
	main_handler = handlers.RateLimitRequest(main_handler)
	main_handler = handlers.LogRequest(main_handler)

	// run sass when in development mode
	if config.GlobalConfig.InDevelopment {
		go func() {
			cmd := exec.Command("sass", "--watch", "assets/css")
			cmd.Stdout = os.Stderr // Sass quirk
			cmd.Stderr = os.Stderr
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pdeathsig: syscall.SIGHUP}
			runtime.LockOSThread() // Go quirk https://github.com/golang/go/issues/27505
			err := cmd.Run()
			if err != nil {
				log.Print(fmt.Errorf("when running sass: %w", err))
			}
		}()
	}

	// Listen
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
	http.Serve(l, http.HandlerFunc(main_handler))
}
