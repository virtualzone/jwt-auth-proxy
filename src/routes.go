package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v10"

	"github.com/gorilla/mux"
)

type Route interface {
	setupRoutes(s *mux.Router)
}

func SendNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func SendBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func SendUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func SendAleadyExists(w http.ResponseWriter) {
	w.WriteHeader(http.StatusConflict)
}

func SendCreated(w http.ResponseWriter, id primitive.ObjectID) {
	w.Header().Set("X-Object-ID", id.Hex())
	w.WriteHeader(http.StatusCreated)
}

func SendUpdated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func SendInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func SendJSON(w http.ResponseWriter, v interface{}) {
	json, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func UnmarshalBody(r *http.Request, o interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, &o); err != nil {
		return err
	}
	return nil
}

func UnmarshalValidateBody(r *http.Request, o interface{}) error {
	err := UnmarshalBody(r, &o)
	if err != nil {
		return err
	}
	v := validator.New()
	err = v.Struct(o)
	if err != nil {
		return err
	}
	return nil
}

func GetUserIDFromContext(r *http.Request) string {
	userID := r.Context().Value("UserID")
	if userID == nil {
		return ""
	}
	return userID.(string)
}

func SetCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", GetConfig().CorsOrigin)
	w.Header().Set("Access-Control-Allow-Headers", GetConfig().CorsHeaders)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetCorsHeaders(w)
		next.ServeHTTP(w, r)
	})
}

func VerifyJwtMiddleware(next http.Handler) http.Handler {
	var IsWhitelisted = func(r *http.Request) bool {
		url := r.URL.RequestURI()
		for _, whitelistedURL := range unauthorizedRoutes {
			if strings.HasPrefix(url, whitelistedURL) {
				return true
			}
		}
		for _, whitelistedURL := range GetConfig().ProxyWhitelist {
			if strings.HasPrefix(url, whitelistedURL) {
				return true
			}
		}
		return false
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsWhitelisted(r) {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("JWT header verification failed: missing auth header")
			SendUnauthorized(w)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("JWT header verification failed: invalid auth header")
			SendUnauthorized(w)
			return
		}
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(GetConfig().JwtSigningKey), nil
		})
		if err != nil {
			log.Println("JWT header verification failed: parsing JWT failed with", err)
			SendUnauthorized(w)
			return
		}
		if !token.Valid {
			log.Println("JWT header verification failed: invalid JWT")
			SendUnauthorized(w)
			return
		}
		log.Println("Successfully verified JWT header for UserID", claims.UserID)
		ctx := context.WithValue(r.Context(), "UserID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CorsHandler(w http.ResponseWriter, r *http.Request) {
	SetCorsHeaders(w)
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	var getScheme = func(s string) string {
		if r.URL.Scheme == "" {
			return "http"
		}
		return r.URL.Scheme
	}

	url := r.URL.RequestURI()
	log.Println("Proxying request for", url)

	r.Header.Set("X-Forwarded-For", r.RemoteAddr)
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("X-Forwarded-Proto", getScheme(r.URL.Scheme))
	r.Header.Set("Forwarded", fmt.Sprintf("for=%s;host=%s;proto=%s", r.RemoteAddr, r.Host, getScheme(r.URL.Scheme)))
	r.Header.Set("X-Auth-UserID", GetUserIDFromContext(r))

	target := GetConfig().ProxyTarget
	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.Host = target.Host

	GetApp().Proxy.ServeHTTP(w, r)
}

var unauthorizedRoutes = [...]string{
	GetConfig().PublicAPIPath + "login",
	GetConfig().PublicAPIPath + "signup",
	GetConfig().PublicAPIPath + "confirm/",
	GetConfig().PublicAPIPath + "initpwreset",
}
