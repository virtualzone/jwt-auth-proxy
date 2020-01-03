package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateUser(t *testing.T) {
	clearTestDB()

	payload := `{"email": "foo@bar.com", "password": "12345678", "confirmed": true, "enabled": true, "data": {"color": "blue"}}`
	req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)
	userID := res.Header().Get("X-Object-Id")

	user := GetUserRepository().GetOne(userID)
	if user == nil {
		t.Fatal("Expected user not to be nil")
	}
	if !GetUserRepository().CheckPassword(user.HashedPassword, "12345678") {
		t.Error("Expected hashed password to match")
	}

	req, _ = http.NewRequest("GET", "/users/"+userID, nil)
	res = executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	var data dummyUser
	if err = json.Unmarshal(body, &data); err != nil {
		t.Error(err)
	}

	checkTestString(t, "foo@bar.com", data.Email)
	if !user.Enabled {
		t.Error("Expected user to be enabled")
	}
	if !data.Confirmed {
		t.Error("Expected user to be confirmed")
	}
	checkTestString(t, "blue", data.Data.Color)
}

func TestCreateUserTwice(t *testing.T) {
	clearTestDB()

	payload := `{"email": "foo@bar.com", "password": "12345678", "confirmed": true, "enabled": true, "data": {"color": "blue"}}`
	req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusCreated, res.Code)

	payload = `{"email": "fOo@bAr.com", "password": "12345678", "confirmed": true, "enabled": true, "data": {"color": "blue"}}`
	req, _ = http.NewRequest("POST", "/users/", bytes.NewBufferString(payload))
	res = executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestCreateUserInvalidEmail(t *testing.T) {
	clearTestDB()

	payload := `{"email": "foobar.com", "password": "12345678", "confirmed": true, "enabled": true, "data": {"color": "blue"}}`
	req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestCreateUserShortPassword(t *testing.T) {
	clearTestDB()

	payload := `{"email": "foo@bar.com", "password": "1234567", "confirmed": true, "enabled": true, "data": {"color": "blue"}}`
	req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestCreateUserConflictingPendingChange(t *testing.T) {
	clearTestDB()
	pa := PendingAction{
		ActionType: PendingActionTypeChangeEmail,
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * GetConfig().PendingActionLifetime),
		UserID:     primitive.NewObjectID(),
		Payload:    "foo@bar.com",
		Token:      GetPendingActionRepository().FindUnusedToken(),
	}
	GetPendingActionRepository().Create(&pa)

	payload := `{"email": "fOo@Bar.com", "password": "12345678", "confirmed": true, "enabled": true, "data": {"color": "blue"}}`
	req, _ := http.NewRequest("POST", "/users/", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestGetSimpleUser(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	req, _ := http.NewRequest("GET", "/users/"+user.ID.Hex(), nil)
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	var data dummyUser
	if err = json.Unmarshal(body, &data); err != nil {
		t.Error(err)
	}
	checkTestString(t, user.Email, data.Email)
	checkTestString(t, user.ID.Hex(), data.ID)
	if data.Data != (dummyUserData{}) {
		t.Error("Expected empty user data")
	}
}

func TestGetNonExistingUser(t *testing.T) {
	clearTestDB()
	createTestUser(true)

	req, _ := http.NewRequest("GET", "/users/123456789", nil)
	res := executeBackendTestRequest(req)

	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestDeleteUser(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	req, _ := http.NewRequest("DELETE", "/users/"+user.ID.Hex(), nil)
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req, _ = http.NewRequest("GET", "/users/"+user.ID.Hex(), nil)
	res = executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestDeleteNonExistingUser(t *testing.T) {
	clearTestDB()
	createTestUser(true)

	req, _ := http.NewRequest("DELETE", "/users/123456789", nil)
	res := executeBackendTestRequest(req)

	checkTestResponseCode(t, http.StatusNotFound, res.Code)
}

func TestSetPassword(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	payload := `{"password": "x1x2x3x4"}`
	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/password", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	payload = `{"email": "foo@bar.com", "password": "12345678"}`
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(payload))
	res = executePublicTestRequest(req)
	checkTestResponseCode(t, http.StatusUnauthorized, res.Code)

	payload = `{"email": "foo@bar.com", "password": "x1x2x3x4"}`
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(payload))
	res = executePublicTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
}

func TestSetPasswordShort(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	payload := `{"password": "x1x2x3"}`
	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/password", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestChangeEmail(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	payload := `{"email": "foo2@bar.com"}`
	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/email", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	payload = `{"email": "foo@bar.com", "password": "12345678"}`
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(payload))
	res = executePublicTestRequest(req)
	checkTestResponseCode(t, http.StatusUnauthorized, res.Code)

	payload = `{"email": "foo2@bar.com", "password": "12345678"}`
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(payload))
	res = executePublicTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
}

func TestChangeEmailConflictingChange(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)
	pa := PendingAction{
		ActionType: PendingActionTypeChangeEmail,
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * GetConfig().PendingActionLifetime),
		UserID:     primitive.NewObjectID(),
		Payload:    "foo2@bar.com",
		Token:      GetPendingActionRepository().FindUnusedToken(),
	}
	GetPendingActionRepository().Create(&pa)

	payload := `{"email": "fOo2@bAr.com"}`
	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/email", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestChangeEmailAlreadyExists(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)
	user2 := &User{
		Email: "foo2@bar.com",
	}
	GetUserRepository().Create(user2)

	payload := `{"email": "fOo2@bAr.com"}`
	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/email", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusConflict, res.Code)
}

func TestChangeEmailInvalidAddress(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	payload := `{"email": "foobar.com"}`
	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/email", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestDisableEnableUser(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/disable", nil)
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	payload := `{"email": "foo@bar.com", "password": "12345678"}`
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(payload))
	res = executePublicTestRequest(req)
	checkTestResponseCode(t, http.StatusUnauthorized, res.Code)

	req, _ = http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/enable", nil)
	res = executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	payload = `{"email": "foo@bar.com", "password": "12345678"}`
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(payload))
	res = executePublicTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
}

func TestSetUserData(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	payload := `{"color": "red", "height": 1.85}`
	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/data", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)
}

func TestGetUserData(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	payload := `{"color": "red", "height": 1.85}`
	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/data", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req, _ = http.NewRequest("GET", "/users/"+user.ID.Hex()+"/data", nil)
	res = executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	var data dummyUserData
	if err = json.Unmarshal(body, &data); err != nil {
		t.Error(err)
	}
	checkTestString(t, "red", data.Color)
	if data.Height != 1.85 {
		t.Errorf("Expected 1.85, got %f", data.Height)
	}
}

func TestGetUserWithData(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	payload := `{"color": "red", "height": 1.85}`
	req, _ := http.NewRequest("PUT", "/users/"+user.ID.Hex()+"/data", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusNoContent, res.Code)

	req, _ = http.NewRequest("GET", "/users/"+user.ID.Hex(), nil)
	res = executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	var data dummyUser
	if err = json.Unmarshal(body, &data); err != nil {
		t.Error(err)
	}
	checkTestString(t, user.Email, data.Email)
	checkTestString(t, user.ID.Hex(), data.ID)
	checkTestString(t, "red", data.Data.Color)
	if data.Data.Height != 1.85 {
		t.Errorf("Expected 1.85, got %f", data.Data.Height)
	}
}

func TestCheckPasswordPositive(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	payload := `{"password": "12345678"}`
	req, _ := http.NewRequest("POST", "/users/"+user.ID.Hex()+"/checkpw", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	var result BoolResult
	json.Unmarshal(res.Body.Bytes(), &result)
	if !result.Result {
		t.Error("Expected positive result")
	}
}

func TestCheckPasswordNegative(t *testing.T) {
	clearTestDB()
	user := createTestUser(true)

	payload := `{"password": "00000000"}`
	req, _ := http.NewRequest("POST", "/users/"+user.ID.Hex()+"/checkpw", bytes.NewBufferString(payload))
	res := executeBackendTestRequest(req)
	checkTestResponseCode(t, http.StatusOK, res.Code)

	var result BoolResult
	json.Unmarshal(res.Body.Bytes(), &result)
	if result.Result {
		t.Error("Expected negative result")
	}
}

type dummyUser struct {
	ID         string        `json:"id"`
	Email      string        `json:"email"`
	Confirmed  bool          `json:"confirmed"`
	Enabled    bool          `json:"enabled"`
	CreateDate time.Time     `json:"createDate"`
	Data       dummyUserData `json:"data,omitempty"`
}

type dummyUserData struct {
	Color  string  `json:"color"`
	Height float32 `json:"height"`
}
