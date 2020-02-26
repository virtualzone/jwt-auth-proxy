package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var staticFilePath = getEnv("STATIC_FILE_PATH", "../frontend/build/")
var proxyAddr = getEnv("PROXY_ADDR", "https://proxy:8443/")

type HttpHandler struct {
}

type UserData struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Confirmed  bool      `json:"confirmed"`
	Enabled    bool      `json:"enabled"`
	CreateDate time.Time `json:"createDate"`
}

func (h *HttpHandler) createHttpClient() *http.Client {
	clientCert := getEnv("CLIENT_CERT", "/app/certs/client.crt")
	clientKey := getEnv("CLIENT_KEY", "/app/certs/client.key")
	caCertFile := getEnv("CA_CERT", "/app/certs/ca.crt")
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatal(err)
	}
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{cert},
		},
	}
	client := &http.Client{Transport: tr}
	return client
}

func (h *HttpHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticFilePath+"index.html")
}

func (h *HttpHandler) UserInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Auth-UserID")
	if userID == "" {
		log.Println("Empty X-Auth-UserID request header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Incoming UserInfo request for UserID ", userID)
	client := h.createHttpClient()
	resp, err := client.Get(proxyAddr + "users/" + userID)
	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Error finding user ", userID, ", got HTTP Status Code ", resp.StatusCode)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var data UserData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Println("Error on Unmarshal: ", err)
	}
	json, err := json.Marshal(&data)
	if err != nil {
		log.Println("Error on Marshal: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func prepareHttpServer(server *http.Server) {
	handler := &HttpHandler{}
	http.HandleFunc("/api/userinfo", handler.UserInfo)
	fs := http.FileServer(http.Dir(staticFilePath))
	http.Handle("/", fs)
	http.HandleFunc("/login.html", handler.ServeIndex)
	http.HandleFunc("/signup.html", handler.ServeIndex)
	http.HandleFunc("/confirm.html", handler.ServeIndex)
	http.HandleFunc("/dashboard.html", handler.ServeIndex)
}

func getEnv(key, defaultValue string) string {
	res := os.Getenv(key)
	if res == "" {
		return defaultValue
	}
	return res
}

func main() {
	log.Println("Starting Example Backend Server...")
	var server *http.Server = &http.Server{
		Addr: "0.0.0.0:8090",
	}
	prepareHttpServer(server)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}()
	log.Println("Example Backend HTTP Server listening on 0.0.0.0:8090")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	server.Shutdown(ctx)
}
