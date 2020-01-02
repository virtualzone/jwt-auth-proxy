package main

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestProxyUnauthorized(t *testing.T) {
	clearTestDB()
	createTestUser(true)

	req := newHTTPRequest("GET", "/some/route/test.html", "", nil)
	res := executePublicTestRequest(req)
	checkTestResponseCode(t, http.StatusUnauthorized, res.Code)
}

func TestProxyTargetDown(t *testing.T) {
	clearTestDB()
	loginResponse := createLoginTestUser()

	req := newHTTPRequest("GET", "/some/route/test.html", loginResponse.AccessToken, nil)
	res := executePublicTestRequest(req)
	checkTestResponseCode(t, http.StatusBadGateway, res.Code)
}

func TestProxySuccessWithAuth(t *testing.T) {
	handler := &dummyProxyHandler{}
	var proxy *http.Server = &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: handler,
	}
	go func() {
		proxy.ListenAndServe()
	}()

	clearTestDB()
	user := createTestUser(true)
	loginResponse := loginUser("foo@bar.com", "12345678")

	req := newHTTPRequest("GET", "/some/route/test.html", loginResponse.AccessToken, nil)
	res := executePublicTestRequest(req)

	proxy.Shutdown(context.TODO())
	checkTestResponseCode(t, http.StatusOK, res.Code)
	if handler.Headers.Get("X-Auth-UserID") != user.ID.Hex() {
		t.Error("Expected X-Auth-UserID header to match actual User ID")
	}
	if !strings.HasPrefix(handler.Headers.Get("Authorization"), "Bearer ") {
		t.Error("Expected Authorization: Bearer [...] header")
	}
}

func TestProxySuccessWhitelisted(t *testing.T) {
	handler := &dummyProxyHandler{}
	var proxy *http.Server = &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: handler,
	}
	go func() {
		proxy.ListenAndServe()
	}()

	clearTestDB()

	req := newHTTPRequest("GET", "/some/route/whitelist.html", "", nil)
	res := executePublicTestRequest(req)

	proxy.Shutdown(context.TODO())
	checkTestResponseCode(t, http.StatusOK, res.Code)
	if handler.Headers.Get("X-Auth-UserID") != "" {
		t.Error("Expected empty X-Auth-UserID header")
	}
	if handler.Headers.Get("Authorization") != "" {
		t.Error("Expected empty Authorizationheader")
	}
}

type dummyProxyHandler struct {
	Headers http.Header
}

func (h *dummyProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Headers = r.Header
}
