package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMain(m *testing.M) {
	os.Setenv("PROXY_TARGET", "http://127.0.0.1:8090")
	os.Setenv("PROXY_WHITELIST", "/some/route/whitelist.html:/some/whitelist")
	os.Setenv("TEMPLATE_SIGNUP", "../test/res/signup.tpl")
	os.Setenv("TEMPLATE_CHANGE_EMAIL", "../test/res/changeemail.tpl")
	os.Setenv("TEMPLATE_RESET_PASSWORD", "../test/res/resetpassword.tpl")
	os.Setenv("TEMPLATE_NEW_PASSWORD", "../test/res/newpassword.tpl")
	GetConfig().ReadConfig()
	smtpClient = func(addr string) (dialer, error) {
		client := &smtpDialerMock{}
		return client, nil
	}
	a := GetApp()
	GetDatatabase().connectMongoDb("mongodb://localhost:27017", "jwt_auth_proxy_test")
	a.InitializePublicRouter()
	a.InitializeBackendRouter()
	readMailTemplatesFromFile()
	code := m.Run()
	GetDatatabase().disconnect()
	os.Exit(code)
}

func createTestUser(confirmed bool) *User {
	user := &User{
		Email:          "foo@bar.com",
		CreateDate:     time.Now(),
		HashedPassword: GetUserRepository().GetHashedPassword("12345678"),
		Confirmed:      confirmed,
		Enabled:        true,
	}
	GetUserRepository().Create(user)
	return user
}

func loginUser(username, password string) *LoginResponse {
	payload := "{\"email\": \"" + username + "\", \"password\": \"" + password + "\"}"
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(payload))
	res := executePublicTestRequest(req)
	var loginResponse LoginResponse
	json.Unmarshal(res.Body.Bytes(), &loginResponse)
	return &loginResponse
}

func createLoginTestUser() *LoginResponse {
	createTestUser(true)
	return loginUser("foo@bar.com", "12345678")
}

func newHTTPRequest(method, url, accessToken string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	return req
}

func clearTestDB() {
	GetPendingActionRepository().GetCollection().DeleteMany(context.TODO(), bson.D{})
	GetRefreshTokenRepository().GetCollection().DeleteMany(context.TODO(), bson.D{})
	GetUserRepository().GetCollection().DeleteMany(context.TODO(), bson.D{})
}

func executePublicTestRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	GetApp().PublicRouter.ServeHTTP(rr, req)
	return rr
}

func executeBackendTestRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	GetApp().BackendRouter.ServeHTTP(rr, req)
	return rr
}

func checkTestResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Expected HTTP Status %d, but got %d", expected, actual)
	}
}

func checkTestString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Expected '%s', but got '%s'", expected, actual)
	}
}
