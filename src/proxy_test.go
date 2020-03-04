package main

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestProxyUnauthorizedNoMatch(t *testing.T) {
	clearTestDB()
	createTestUser(true)

	req := newHTTPRequest("GET", "/some/route/test.html", "", nil)
	res := executePublicTestRequest(req)
	checkTestResponseCode(t, http.StatusUnauthorized, res.Code)
}

func TestProxyUnauthorizedPrefixMatch(t *testing.T) {
	clearTestDB()
	createTestUser(true)

	req := newHTTPRequest("GET", "/some/whitelist2", "", nil)
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

func TestProxySuccessWhitelistWithAuth(t *testing.T) {
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

	req := newHTTPRequest("GET", "/some/whitelist/page.html", loginResponse.AccessToken, nil)
	res := executePublicTestRequest(req)

	proxy.Shutdown(context.TODO())
	checkTestResponseCode(t, http.StatusOK, res.Code)
	if !strings.HasPrefix(handler.Headers.Get("Authorization"), "Bearer ") {
		t.Error("Expected Authorization: Bearer [...] header")
	}
	if handler.Headers.Get("X-Auth-UserID") != user.ID.Hex() {
		t.Error("Expected X-Auth-UserID header to match actual User ID '" + user.ID.Hex() + "' but got '" + handler.Headers.Get("X-Auth-UserID") + "'")
	}
}

func TestProxyFakeUserIDHeaderWhitelistWithoutAuth(t *testing.T) {
	handler := &dummyProxyHandler{}
	var proxy *http.Server = &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: handler,
	}
	go func() {
		proxy.ListenAndServe()
	}()

	clearTestDB()

	req := newHTTPRequest("GET", "/some/whitelist/page.html", "", nil)
	req.Header.Set("X-Auth-UserID", "FAKE")
	res := executePublicTestRequest(req)

	proxy.Shutdown(context.TODO())
	checkTestResponseCode(t, http.StatusOK, res.Code)
	if handler.Headers.Get("Authorization") != "" {
		t.Error("Expected empty Authorizationheader, got: " + handler.Headers.Get("Authorization"))
	}
	if handler.Headers.Get("X-Auth-UserID") != "" {
		t.Error("Expected empty X-Auth-UserID header, got: " + handler.Headers.Get("X-Auth-UserID"))
	}
}

func TestProxyFakeAuthHeaderWhitelistWithoutAuth(t *testing.T) {
	handler := &dummyProxyHandler{}
	var proxy *http.Server = &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: handler,
	}
	go func() {
		proxy.ListenAndServe()
	}()

	clearTestDB()

	req := newHTTPRequest("GET", "/some/whitelist/page.html", "", nil)
	req.Header.Set("Authorization", "Bearer FAKE")
	res := executePublicTestRequest(req)

	proxy.Shutdown(context.TODO())
	checkTestResponseCode(t, http.StatusOK, res.Code)
	if handler.Headers.Get("Authorization") != "" {
		t.Error("Expected empty Authorizationheader, got: " + handler.Headers.Get("Authorization"))
	}
	if handler.Headers.Get("X-Auth-UserID") != "" {
		t.Error("Expected empty X-Auth-UserID header, got: " + handler.Headers.Get("X-Auth-UserID"))
	}
}

func TestProxySuccessWhitelistedWithoutAuth(t *testing.T) {
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

func TestProxySuccessWhitelistedSubPath(t *testing.T) {
	handler := &dummyProxyHandler{}
	var proxy *http.Server = &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: handler,
	}
	go func() {
		proxy.ListenAndServe()
	}()

	clearTestDB()

	req := newHTTPRequest("GET", "/some/whitelist/test.html", "", nil)
	res := executePublicTestRequest(req)

	proxy.Shutdown(context.TODO())
	checkTestResponseCode(t, http.StatusOK, res.Code)
}

type dummyProxyHandler struct {
	Headers http.Header
}

func (h *dummyProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Headers = r.Header
}
