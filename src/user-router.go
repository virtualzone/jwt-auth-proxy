package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type UserRouter struct {
}

func (router *UserRouter) setupRoutes(s *mux.Router) {
	s.HandleFunc("/{id}", router.getOne).Methods("GET")
	s.HandleFunc("/{id}", router.delete).Methods("DELETE")
	s.HandleFunc("/{id}/email", router.setEmail).Methods("PUT")
	s.HandleFunc("/{id}/password", router.setPassword).Methods("PUT")
	s.HandleFunc("/{id}/enable", router.enableUser).Methods("PUT")
	s.HandleFunc("/{id}/disable", router.disableUser).Methods("PUT")
	s.HandleFunc("/{id}/data", router.getUserData).Methods("GET")
	s.HandleFunc("/{id}/data", router.setUserData).Methods("PUT")
	s.HandleFunc("/{id}/checkpw", router.checkPassword).Methods("POST")
	s.HandleFunc("/", router.Create).Methods("POST")
	s.HandleFunc("/", router.getAll).Methods("GET")
}

func (router *UserRouter) Create(w http.ResponseWriter, r *http.Request) {
	var data CreateUserRequest
	if UnmarshalValidateBody(r, &data) != nil {
		log.Println("Received invalid create user request")
		SendBadRequest(w)
		return
	}
	user := &User{
		Email:          data.Email,
		HashedPassword: GetUserRepository().GetHashedPassword(data.Password),
		Confirmed:      data.Confirmed,
		Enabled:        data.Enabled,
		Data:           data.Data,
		CreateDate:     time.Now(),
	}
	GetUserRepository().Create(user)
	SendCreated(w, user.ID)
}

func (router *UserRouter) getOne(w http.ResponseWriter, r *http.Request) {
	user := router.getUserFromMuxVars(w, r)
	if user == nil {
		SendNotFound(w)
		return
	}
	user.HashedPassword = ""
	data, err := router.prepareUserData(user)
	if err != nil {
		SendInternalServerError(w)
		return
	}
	user.Data = data
	SendJSON(w, user)
}

func (router *UserRouter) delete(w http.ResponseWriter, r *http.Request) {
	user := router.getUserFromMuxVars(w, r)
	if user == nil {
		SendNotFound(w)
		return
	}
	GetUserRepository().Delete(user)
	SendUpdated(w)
}

func (router *UserRouter) setEmail(w http.ResponseWriter, r *http.Request) {
	user := router.getUserFromMuxVars(w, r)
	if user == nil {
		SendNotFound(w)
		return
	}
	var data SetEmailRequest
	if UnmarshalValidateBody(r, &data) != nil {
		SendBadRequest(w)
		return
	}
	user.Email = data.Email
	GetUserRepository().Update(user)
	SendUpdated(w)
}

func (router *UserRouter) setPassword(w http.ResponseWriter, r *http.Request) {
	user := router.getUserFromMuxVars(w, r)
	if user == nil {
		SendNotFound(w)
		return
	}
	var data SetPasswordRequest
	if UnmarshalValidateBody(r, &data) != nil {
		SendBadRequest(w)
		return
	}
	user.HashedPassword = GetUserRepository().GetHashedPassword(data.Password)
	GetUserRepository().Update(user)
	SendUpdated(w)
}

func (router *UserRouter) disableUser(w http.ResponseWriter, r *http.Request) {
	user := router.getUserFromMuxVars(w, r)
	if user == nil {
		SendNotFound(w)
		return
	}
	user.Enabled = false
	GetUserRepository().Update(user)
	SendUpdated(w)
}

func (router *UserRouter) enableUser(w http.ResponseWriter, r *http.Request) {
	user := router.getUserFromMuxVars(w, r)
	if user == nil {
		SendNotFound(w)
		return
	}
	user.Enabled = true
	GetUserRepository().Update(user)
	SendUpdated(w)
}

func (router *UserRouter) setUserData(w http.ResponseWriter, r *http.Request) {
	user := router.getUserFromMuxVars(w, r)
	if user == nil {
		SendNotFound(w)
		return
	}
	var data interface{}
	if err := UnmarshalBody(r, &data); err != nil {
		SendBadRequest(w)
		return
	}
	user.Data = data
	GetUserRepository().Update(user)
	SendUpdated(w)
}

func (router *UserRouter) getUserData(w http.ResponseWriter, r *http.Request) {
	user := router.getUserFromMuxVars(w, r)
	if user == nil {
		SendNotFound(w)
		return
	}
	data, err := router.prepareUserData(user)
	if err != nil {
		SendInternalServerError(w)
		return
	}
	SendJSON(w, data)
}

func (router *UserRouter) checkPassword(w http.ResponseWriter, r *http.Request) {
	user := router.getUserFromMuxVars(w, r)
	if user == nil {
		SendNotFound(w)
		return
	}
	var data SetPasswordRequest
	if UnmarshalValidateBody(r, &data) != nil {
		SendBadRequest(w)
		return
	}
	result := &BoolResult{
		Result: GetUserRepository().CheckPassword(user.HashedPassword, data.Password),
	}
	SendJSON(w, result)
}

func (router *UserRouter) getAll(w http.ResponseWriter, r *http.Request) {
	// TODO Implement method
	SendInternalServerError(w)
}

func (router *UserRouter) getUserFromMuxVars(w http.ResponseWriter, r *http.Request) *User {
	vars := mux.Vars(r)
	user := GetUserRepository().GetOne(vars["id"])
	if user == nil {
		return nil
	}
	return user
}

func (router *UserRouter) prepareUserData(user *User) (map[string]interface{}, error) {
	m, err := json.Marshal(user.Data)
	if err != nil {
		return nil, err
	}
	var data []keyValue
	if err := json.Unmarshal(m, &data); err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	for i := 0; i < len(data); i++ {
		item := data[i]
		res[item.Key] = item.Value
	}
	return res, nil
}

type SetEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type SetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type BoolResult struct {
	Result bool `json:"result"`
}

type keyValue struct {
	Key   string
	Value interface{}
}

type CreateUserRequest struct {
	Email     string      `json:"email" validate:"required,email"`
	Password  string      `json:"password" validate:"required,min=8,max=32"`
	Confirmed bool        `json:"confirmed,omitempty"`
	Enabled   bool        `json:"enabled,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}
