package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var _appInstance *App
var _appOnce sync.Once

func GetApp() *App {
	_appOnce.Do(func() {
		_appInstance = &App{}
	})
	return _appInstance
}

type App struct {
	PublicRouter              *mux.Router
	BackendRouter             *mux.Router
	Proxy                     *httputil.ReverseProxy
	CleanRefreshTokensTicker  *time.Ticker
	CleanPendingActionsTicker *time.Ticker
}

func (a *App) InitializePublicRouter() {
	a.InitializeProxy()
	a.PublicRouter = mux.NewRouter()
	routers := make(map[string]Route)
	routers[GetConfig().PublicAPIPath] = &AuthRouter{}
	for route, router := range routers {
		subRouter := a.PublicRouter.PathPrefix(route).Subrouter()
		router.setupRoutes(subRouter)
	}
	if GetConfig().EnableCors {
		a.PublicRouter.PathPrefix("/").Methods("OPTIONS").HandlerFunc(CorsHandler)
		a.PublicRouter.Use(CorsMiddleware)
	}
	a.PublicRouter.PathPrefix("/").HandlerFunc(ProxyHandler)
	a.PublicRouter.Use(VerifyJwtMiddleware)
}

func (a *App) InitializeBackendRouter() {
	a.BackendRouter = mux.NewRouter()
	routers := make(map[string]Route)
	routers["/users/"] = &UserRouter{}
	for route, router := range routers {
		subRouter := a.BackendRouter.PathPrefix(route).Subrouter()
		router.setupRoutes(subRouter)
	}
}

func (a *App) InitializeProxy() {
	target := GetConfig().ProxyTarget
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = a._SingleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	a.Proxy = &httputil.ReverseProxy{Director: director}
}

func (a *App) InitializeTimers() {
	a.CleanRefreshTokensTicker = time.NewTicker(time.Hour * 1)
	go func() {
		for {
			select {
			case <-a.CleanRefreshTokensTicker.C:
				log.Println("Cleaning up expired refresh tokens...")
				GetRefreshTokenRepository().CleanUp()
			}
		}
	}()
	a.CleanPendingActionsTicker = time.NewTicker(time.Hour * 1)
	go func() {
		for {
			select {
			case <-a.CleanPendingActionsTicker.C:
				log.Println("Cleaning up expired pending actions...")
				GetPendingActionRepository().CleanUp()
			}
		}
	}()
}

func (a *App) GenerateBackendCert() {
	log.Println("Generating Backend mTLS Certificate...")
	dir := GetConfig().BackendCertDir
	ca, err := CertCreateCA()
	if err != nil {
		log.Fatalln(err)
	}
	if err := ca.SavePrivateKey(dir + "ca.key"); err != nil {
		log.Fatalln(err)
	}
	if err := ca.SaveCertificate(dir + "ca.crt"); err != nil {
		log.Fatalln(err)
	}

	server, err := CertCreateSign(ca)
	if err != nil {
		log.Fatalln(err)
	}
	if err := server.SavePrivateKey(dir + "server.key"); err != nil {
		log.Fatalln(err)
	}
	if err := server.SaveCertificate(dir + "server.crt"); err != nil {
		log.Fatalln(err)
	}

	client, err := CertCreateSign(ca)
	if err != nil {
		log.Fatalln(err)
	}
	if err := client.SavePrivateKey(dir + "client.key"); err != nil {
		log.Fatalln(err)
	}
	if err := client.SaveCertificate(dir + "client.crt"); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) Run(publicListenAddr, backendListenAddr string) {
	if GetConfig().BackendGenerateCert {
		a.GenerateBackendCert()
	}
	log.Println("Initializing REST services...")
	publicServer := &http.Server{
		Addr:         publicListenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.PublicRouter,
	}
	go func() {
		if err := publicServer.ListenAndServe(); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}()
	log.Println("Public HTTP Server listening on", publicListenAddr)
	tlsConfig := a._CreateTLSConfig()
	backendServer := &http.Server{
		Addr:         backendListenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.BackendRouter,
		TLSConfig:    tlsConfig,
	}
	go func() {
		if err := backendServer.ListenAndServeTLS(GetConfig().BackendCertDir+"server.crt", GetConfig().BackendCertDir+"server.key"); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}()
	log.Println("Backend HTTPS Server listening on", backendListenAddr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	a.CleanPendingActionsTicker.Stop()
	a.CleanRefreshTokensTicker.Stop()
	backendServer.Shutdown(ctx)
	publicServer.Shutdown(ctx)
}

func (a *App) _CreateTLSConfig() *tls.Config {
	crt, _ := filepath.Abs(GetConfig().BackendCertDir + "ca.crt")
	clientCaCert, err := ioutil.ReadFile(crt)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(clientCaCert)
	tlsConfig := &tls.Config{
		ClientCAs:                caCertPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig
}

func (app *App) _SingleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
