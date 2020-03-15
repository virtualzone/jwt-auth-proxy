package main

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

// AuthRouter handles authentication related REST requests
type AuthRouter struct {
}

func (router *AuthRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/login", router.Login).Methods("POST")
	s.HandleFunc("/refresh", router.Refresh).Methods("POST")
	s.HandleFunc("/logout", router.Logout).Methods("POST")
	s.HandleFunc("/ping", router.Ping).Methods("GET")
	if GetConfig().AllowSignup {
		s.HandleFunc("/signup", router.Signup).Methods("POST")
	}
	if GetConfig().AllowChangePassword {
		s.HandleFunc("/setpw", router.ChangePassword).Methods("POST")
	}
	if GetConfig().AllowChangeEmail {
		s.HandleFunc("/changeemail", router.ChangeEmail).Methods("POST")
	}
	if GetConfig().AllowForgotPassword {
		s.HandleFunc("/initpwreset", router.InitForgotPassword).Methods("POST")
	}
	if GetConfig().AllowDeleteAccount {
		s.HandleFunc("/delete", router.DeleteAccount).Methods("POST")
	}
	if GetConfig().EnableTOTP {
		s.HandleFunc("/otp/init", router.OTPInit).Methods("POST")
		s.HandleFunc("/otp/confirm", router.OTPConfirm).Methods("POST")
		s.HandleFunc("/otp/disable", router.OTPDisable).Methods("POST")
	}
	s.HandleFunc("/confirm/{id}", router.Confirm).Methods("POST")
	s.PathPrefix("/").Methods("OPTIONS").HandlerFunc(CorsHandler)
	s.PathPrefix("/").HandlerFunc(router.NotFound)
}

// NotFound handles all other requests
func (router *AuthRouter) NotFound(w http.ResponseWriter, r *http.Request) {
	SendNotFound(w)
}

// Login handles /login requests
func (router *AuthRouter) Login(w http.ResponseWriter, r *http.Request) {
	var data LoginRequest
	if UnmarshalValidateBody(r, &data) != nil {
		log.Println("Invalid login attempt: failed unmarshalling request")
		SendBadRequest(w)
		return
	}
	user := GetUserRepository().GetByEmail(data.Email)
	if user == nil {
		log.Println("Invalid login attempt: invalid username", data.Email)
		SendUnauthorized(w)
		return
	}
	if user.Confirmed == false {
		log.Println("Invalid login attempt: unconfirmed account", user.ID.Hex())
		SendUnauthorized(w)
		return
	}
	if user.Enabled == false {
		log.Println("Invalid login attempt: disabled account", user.ID.Hex())
		SendUnauthorized(w)
		return
	}
	if GetUserRepository().CheckPassword(user.HashedPassword, data.Password) == false {
		log.Println("Invalid login attempt: invalid password for UserID", user.ID.Hex())
		SendUnauthorized(w)
		return
	}
	if user.OTPEnabled && GetConfig().EnableTOTP {
		if len(strings.TrimSpace(data.OTP)) != 6 {
			log.Println("Login attempt successful, but missing OTP for UserID", user.ID.Hex())
			SendJSON(w, &LoginResponse{RequireOTP: true})
			return
		}
		if !router._IsValidOTP(user, data.OTP) {
			log.Println("Login attempt successful, but OTP invalid for UserID", user.ID.Hex())
			SendJSON(w, &LoginResponse{RequireOTP: true})
			return
		}
	}
	log.Println("Successful login for UserID", user.ID.Hex())
	refreshToken := router._CreateRefreshToken(user)
	accessToken := router._CreateAccessToken(user)
	SendJSON(w, &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})
}

// Refresh handles /refresh requests
func (router *AuthRouter) Refresh(w http.ResponseWriter, r *http.Request) {
	var data RefreshRequest
	if UnmarshalValidateBody(r, &data) != nil {
		log.Println("Invalid token refresh attempt: failed unmarshalling request")
		SendBadRequest(w)
		return
	}
	refreshToken := GetRefreshTokenRepository().GetByToken(data.RefreshToken)
	if refreshToken == nil {
		log.Println("Invalid token refresh attempt: invalid refresh token")
		SendBadRequest(w)
		return
	}
	if refreshToken.ExpiryDate.Before(time.Now()) {
		log.Println("Invalid token refresh attempt: refresh token expired")
		SendBadRequest(w)
		return
	}
	user := GetUserRepository().GetOne(GetUserIDFromContext(r))
	if user == nil {
		log.Println("Invalid token refresh attempt: invalid UserID", GetUserIDFromContext(r))
		SendUnauthorized(w)
		return
	}
	if user.Confirmed == false {
		log.Println("Invalid token refresh attempt: unconfirmed account", user.ID.Hex())
		SendUnauthorized(w)
		return
	}
	if user.Enabled == false {
		log.Println("Invalid token refresh attempt: disabled account", user.ID.Hex())
		SendUnauthorized(w)
		return
	}
	log.Println("Successful token refresh for UserID", user.ID.Hex())
	accessToken := router._CreateAccessToken(user)
	SendJSON(w, &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})
}

// Logout handles /logout requests
func (router *AuthRouter) Logout(w http.ResponseWriter, r *http.Request) {
	var data RefreshRequest
	if UnmarshalValidateBody(r, &data) != nil {
		log.Println("Invalid logout attempt: failed unmarshalling request")
		SendBadRequest(w)
		return
	}
	refreshToken := GetRefreshTokenRepository().GetByToken(data.RefreshToken)
	if refreshToken == nil {
		log.Println("Invalid logout attempt: invalid refresh token")
		SendBadRequest(w)
		return
	}
	GetRefreshTokenRepository().Delete(refreshToken)
	SendUpdated(w)
}

// Ping handles /ping requests
func (router *AuthRouter) Ping(w http.ResponseWriter, r *http.Request) {
	SendUpdated(w)
}

func (router *AuthRouter) _CreateAccessToken(user *User) string {
	claims := &Claims{
		Email:  user.Email,
		UserID: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(GetConfig().AccessTokenLifetime * time.Minute).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	jwtString, err := accessToken.SignedString([]byte(GetConfig().JwtSigningKey))
	if err != nil {
		return ""
	}
	return jwtString
}

// Signup handles /signup requests
func (router *AuthRouter) Signup(w http.ResponseWriter, r *http.Request) {
	var data SignupRequest
	if UnmarshalValidateBody(r, &data) != nil {
		SendBadRequest(w)
		return
	}
	user := GetUserRepository().GetByEmail(data.Email)
	if user != nil {
		SendAleadyExists(w)
		return
	}
	if len(GetPendingActionRepository().GetByPayload(data.Email)) != 0 {
		SendAleadyExists(w)
		return
	}
	user = &User{
		Email:          data.Email,
		HashedPassword: GetUserRepository().GetHashedPassword(data.Password),
		Confirmed:      false,
		Enabled:        true,
		CreateDate:     time.Now(),
	}
	GetUserRepository().Create(user)
	pa := router._CreateConfirmPendingAction(user, PendingActionTypeConfirmAccount, "")
	router._SendWelcomeMailToNewUser(user, pa)
	SendCreated(w, user.ID)
}

// ChangePassword handles /changepw requests
func (router *AuthRouter) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var data ChangePasswordRequest
	if UnmarshalValidateBody(r, &data) != nil {
		log.Println("Invalid change password attempt: failed unmarshalling request")
		SendBadRequest(w)
		return
	}
	user := GetUserRepository().GetOne(GetUserIDFromContext(r))
	if user == nil {
		log.Println("Invalid change password attempt: invalid UserID", GetUserIDFromContext(r))
		SendUnauthorized(w)
		return
	}
	if !GetUserRepository().CheckPassword(user.HashedPassword, data.OldPassword) {
		log.Println("Invalid change password attempt: incorrect old password for UserID", GetUserIDFromContext(r))
		SendUnauthorized(w)
		return
	}
	user.HashedPassword = GetUserRepository().GetHashedPassword(data.NewPassword)
	GetUserRepository().Update(user)
	SendUpdated(w)
}

// ChangeEmail handles /changeemail requests
func (router *AuthRouter) ChangeEmail(w http.ResponseWriter, r *http.Request) {
	var data LoginRequest
	if UnmarshalValidateBody(r, &data) != nil {
		log.Println("Invalid change email attempt: failed unmarshalling request")
		SendBadRequest(w)
		return
	}
	user := GetUserRepository().GetOne(GetUserIDFromContext(r))
	if user == nil {
		log.Println("Invalid change email attempt: invalid UserID", GetUserIDFromContext(r))
		SendUnauthorized(w)
		return
	}
	if !GetUserRepository().CheckPassword(user.HashedPassword, data.Password) {
		log.Println("Invalid change email attempt: incorrect password for UserID", GetUserIDFromContext(r))
		SendUnauthorized(w)
		return
	}
	if GetUserRepository().GetByEmail(data.Email) != nil {
		SendAleadyExists(w)
		return
	}
	if len(GetPendingActionRepository().GetByPayload(data.Email)) != 0 {
		SendAleadyExists(w)
		return
	}
	pa := router._CreateConfirmPendingAction(user, PendingActionTypeChangeEmail, data.Email)
	router._SendConfirmEmailChangeMail(user, pa)
	SendUpdated(w)
}

// InitForgotPassword handles /initpwreset requests
func (router *AuthRouter) InitForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data ForgotPasswordRequest
	if UnmarshalValidateBody(r, &data) != nil {
		log.Println("Invalid init forgot password attempt: failed unmarshalling request")
		SendBadRequest(w)
		return
	}
	user := GetUserRepository().GetByEmail(data.Email)
	if user == nil {
		log.Println("Invalid init forgot password attempt: invalid email", data.Email)
		SendBadRequest(w)
		return
	}
	pa := router._CreateConfirmPendingAction(user, PendingActionTypeInitPasswordReset, "")
	router._SendConfirmPasswordResetMail(user, pa)
	SendUpdated(w)
}

// DeleteAccount handles /delete requests
func (router *AuthRouter) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var data DeleteAccountRequest
	if UnmarshalValidateBody(r, &data) != nil {
		log.Println("Invalid delete account attempt: failed unmarshalling request")
		SendBadRequest(w)
		return
	}
	user := GetUserRepository().GetOne(GetUserIDFromContext(r))
	if user == nil {
		log.Println("Invalid delete account attempt: invalid UserID", GetUserIDFromContext(r))
		SendUnauthorized(w)
		return
	}
	if !GetUserRepository().CheckPassword(user.HashedPassword, data.Password) {
		log.Println("Invalid delete account attempt: incorrect password for UserID", GetUserIDFromContext(r))
		SendUnauthorized(w)
		return
	}
	GetUserRepository().Delete(user)
	SendUpdated(w)
}

// Confirm handles /confirm requests
func (router *AuthRouter) Confirm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("Requested confirm for ID", vars["id"])
	pa := GetPendingActionRepository().GetByToken(vars["id"])
	if pa == nil {
		SendNotFound(w)
		return
	}
	user := GetUserRepository().GetOne(pa.UserID.Hex())
	if user == nil {
		SendNotFound(w)
		return
	}
	if !user.Enabled {
		SendNotFound(w)
		return
	}
	switch pa.ActionType {
	case PendingActionTypeConfirmAccount:
		router._ConfirmAccountActivation(w, pa, user)
		break
	case PendingActionTypeChangeEmail:
		router._ConfirmEmailChange(w, pa, user)
		break
	case PendingActionTypeInitPasswordReset:
		router._ConfirmPasswordReset(w, pa, user)
		break
	default:
		SendInternalServerError(w)
	}
}

func (router *AuthRouter) OTPInit(w http.ResponseWriter, r *http.Request) {
	user := GetUserRepository().GetOne(GetUserIDFromContext(r))
	if user.OTPEnabled && user.OTPSecret != "" {
		SendBadRequest(w)
		return
	}

	options := totp.GenerateOpts{
		Issuer:      GetConfig().TOTPIssuer,
		AccountName: user.Email,
	}
	key, err := totp.Generate(options)
	if err != nil {
		SendInternalServerError(w)
		return
	}

	secret, err := Encrypt(GetConfig().TOTPSecretEncryptionKey, key.Secret())
	if err != nil {
		log.Println("Could not encrypt TOTP secret:", err)
		SendInternalServerError(w)
		return
	}
	user.OTPSecret = secret
	user.OTPEnabled = false
	GetUserRepository().Update(user)

	img, _ := key.Image(256, 256)
	buf := bytes.NewBuffer([]byte{})
	png.Encode(buf, img)
	res := OTPInitResponse{
		Secret: key.Secret(),
		Image:  base64.StdEncoding.EncodeToString(buf.Bytes()),
	}
	SendJSON(w, res)
}

func (router *AuthRouter) OTPDisable(w http.ResponseWriter, r *http.Request) {
	user := GetUserRepository().GetOne(GetUserIDFromContext(r))
	user.OTPSecret = ""
	user.OTPEnabled = false
	GetUserRepository().Update(user)
	SendUpdated(w)
}

func (router *AuthRouter) OTPConfirm(w http.ResponseWriter, r *http.Request) {
	var data OTPValidateRequest
	if UnmarshalValidateBody(r, &data) != nil {
		SendBadRequest(w)
		return
	}
	user := GetUserRepository().GetOne(GetUserIDFromContext(r))
	if user.OTPEnabled {
		log.Println("Invalid OTP confirm attempt: user has not enabled OTP")
		SendBadRequest(w)
		return
	}
	if strings.TrimSpace(user.OTPSecret) == "" {
		log.Println("Invalid OTP confirm attempt: empty secret")
		SendBadRequest(w)
		return
	}
	if !router._IsValidOTP(user, data.Passcode) {
		log.Println("Invalid OTP confirm attempt: invalid passcode")
		SendBadRequest(w)
		return
	}
	user.OTPEnabled = true
	GetUserRepository().Update(user)
	SendUpdated(w)
}

func (router *AuthRouter) _IsValidOTP(user *User, passcode string) bool {
	secret, err := Decrypt(GetConfig().TOTPSecretEncryptionKey, user.OTPSecret)
	if err != nil {
		log.Println("Could not decrypt TOTP secret:", err)
		return false
	}
	return totp.Validate(passcode, secret)
}

func (router *AuthRouter) _ConfirmAccountActivation(w http.ResponseWriter, pa *PendingAction, user *User) {
	user.Confirmed = true
	GetUserRepository().Update(user)
	GetPendingActionRepository().Delete(pa)
	SendUpdated(w)
}

func (router *AuthRouter) _ConfirmEmailChange(w http.ResponseWriter, pa *PendingAction, user *User) {
	user.Email = pa.Payload
	GetUserRepository().Update(user)
	GetPendingActionRepository().Delete(pa)
	SendUpdated(w)
}

func (router *AuthRouter) _ConfirmPasswordReset(w http.ResponseWriter, pa *PendingAction, user *User) {
	password := GetConfig().GenerateRandomPassword(8)
	user.HashedPassword = GetUserRepository().GetHashedPassword(password)
	GetUserRepository().Update(user)
	GetPendingActionRepository().Delete(pa)
	router._SendNewPassword(user, password)
	SendUpdated(w)
}

func (router *AuthRouter) _CreateRefreshToken(user *User) *RefreshToken {
	e := &RefreshToken{
		Token:      GetRefreshTokenRepository().FindUnusedToken(),
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * GetConfig().RefreshTokenLifetime),
		UserID:     user.ID,
	}
	GetRefreshTokenRepository().Create(e)
	return e
}

func (router *AuthRouter) _CreateConfirmPendingAction(user *User, actionType int, payload string) *PendingAction {
	pa := PendingAction{
		ActionType: actionType,
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * GetConfig().PendingActionLifetime),
		UserID:     user.ID,
		Payload:    payload,
		Token:      GetPendingActionRepository().FindUnusedToken(),
	}
	GetPendingActionRepository().Create(&pa)
	return &pa
}

func (router *AuthRouter) _SendWelcomeMailToNewUser(user *User, pa *PendingAction) {
	var buf bytes.Buffer
	TemplateSignup.Execute(&buf, ConfirmMailVars{
		From:      GetConfig().SMTPSenderAddr,
		To:        user.Email,
		ConfirmID: pa.Token,
	})
	SendMail(user.Email, buf.String())
}

func (router *AuthRouter) _SendConfirmEmailChangeMail(user *User, pa *PendingAction) {
	var buf bytes.Buffer
	TemplateChangeEmail.Execute(&buf, ConfirmMailVars{
		From:      GetConfig().SMTPSenderAddr,
		To:        pa.Payload,
		ConfirmID: pa.Token,
	})
	SendMail(pa.Payload, buf.String())
}

func (router *AuthRouter) _SendConfirmPasswordResetMail(user *User, pa *PendingAction) {
	var buf bytes.Buffer
	TemplateResetPassword.Execute(&buf, ConfirmMailVars{
		From:      GetConfig().SMTPSenderAddr,
		To:        user.Email,
		ConfirmID: pa.Token,
	})
	SendMail(user.Email, buf.String())
}

func (router *AuthRouter) _SendNewPassword(user *User, password string) {
	var buf bytes.Buffer
	TemplateNewPassword.Execute(&buf, PasswordMailVars{
		From:     GetConfig().SMTPSenderAddr,
		To:       user.Email,
		Password: password,
	})
	SendMail(user.Email, buf.String())
}

// LoginRequest holds the POST payload for login requests
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
	OTP      string `json:"otp"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// RefreshRequest holds the POST payload for refresh requests
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// Claims holds payload the issued JWTs
type Claims struct {
	Email  string `json:"email"`
	UserID string `json:"userID"`
	jwt.StandardClaims
}

// LoginResponse holds the response payload for login responses
type LoginResponse struct {
	RequireOTP   bool   `json:"otpRequired"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// ChangePasswordRequest holds the POST payload for password change requests
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,min=8,max=32"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=32"`
}

// SignupRequest holds the POST payload for signup requests
type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

// DeleteAccountRequest holds the POST payload for account delete requests
type DeleteAccountRequest struct {
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type OTPInitResponse struct {
	Secret string `json:"secret"`
	Image  string `json:"image"`
}

type OTPValidateRequest struct {
	Passcode string `json:"passcode" validate:"required,min=6,max=6"`
}
