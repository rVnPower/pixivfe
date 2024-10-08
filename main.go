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
	"sync"
	"syscall"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/server/audit"
	"codeberg.org/vnpower/pixivfe/v2/server/middleware"
	"codeberg.org/vnpower/pixivfe/v2/server/proxy_checker"
	"codeberg.org/vnpower/pixivfe/v2/server/template"
)

func main() {
	if err := config.GlobalConfig.LoadConfig(); err != nil {
		log.Fatalln(err)
	}

	audit.Init(config.GlobalConfig.InDevelopment)
	template.Init(config.GlobalConfig.InDevelopment, "assets/views")

	// Conditionally initialize and start the proxy checker
	if config.GlobalConfig.ProxyCheckEnabled {
		go func() {
			var wg_firstCheck sync.WaitGroup
			wg_firstCheck.Add(1)

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), config.GlobalConfig.ProxyCheckTimeout)
				proxy_checker.CheckProxies(ctx)
				cancel()
				wg_firstCheck.Done()
				if config.GlobalConfig.ProxyCheckInterval != 0 {
					for {
						time.Sleep(config.GlobalConfig.ProxyCheckInterval)
						ctx, cancel := context.WithTimeout(context.Background(), config.GlobalConfig.ProxyCheckTimeout)
						proxy_checker.CheckProxies(ctx)
						cancel()
					}
				}
			}()

			// Wait for the first proxy check to complete
			log.Println("Waiting for initial proxy check to complete...")
			wg_firstCheck.Wait()
			log.Println("Initial proxy check completed.")
		}()
	} else {
		log.Println("Skipping proxy checker initialization.")
	}

	log.Println("Starting server...")

	router := middleware.DefineRoutes()
	// the first middleware is the most outer / first executed one
	router.Use(middleware.ProvideUserContext)  // needed for everything else
	router.Use(middleware.SetLocaleFromCookie) // needed for i18n.*()
	router.Use(middleware.LogRequest)          // all pages need this
	router.Use(middleware.SetPrivacyHeaders)   // all pages need this
	router.Use(middleware.HandleError)         // if the inner handler fails, this shows the error page instead
	router.Use(middleware.InitializeRateLimiter())

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
		log.Print(fmt.Errorf("when running sass: %w", err), err)
	}
}

func chooseListener() net.Listener {
	var listener net.Listener

	// Check if we should use a Unix domain socket
	if config.GlobalConfig.UnixSocket != "" {
		unixAddr := config.GlobalConfig.UnixSocket
		unixListener, err := net.Listen("unix", unixAddr)
		if err != nil {
			// Panic with the error description if unable to listen on Unix socket
			panic(fmt.Errorf("failed to listen on Unix socket %q: %w", unixAddr, err))
		}

		// Assign the listener and log where we are listening
		listener = unixListener
		log.Print(fmt.Sprintf("Listening on Unix domain socket: %v", unixAddr))

	} else {
		// Otherwise, fall back to TCP listener
		addr := net.JoinHostPort(config.GlobalConfig.Host, config.GlobalConfig.Port)
		tcpListener, err := net.Listen("tcp", addr)
		if err != nil {
			// Panic with the error description if unable to listen on TCP
			log.Panic(fmt.Sprintf("Failed to start TCP listener on %v: %v", addr, err))
		}

		// Assign the TCP listener
		listener = tcpListener
		addr = tcpListener.Addr().String()

		// Extract the host and port for logging
		_, port, err := net.SplitHostPort(addr)
		if err != nil {
			// Panic in case of invalid split into host and port
			panic(fmt.Errorf("error parsing listener address %q: %w", addr, err))
		}

		// Log the address and convenient URL for local development
		log.Print(fmt.Sprintf("Listening on %v. Accessible at: http://pixivfe.localhost:%v/", addr, port))
	}

	return listener
}
