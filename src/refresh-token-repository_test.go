package main

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRefreshTokenCleanUp(t *testing.T) {
	clearTestDB()

	t1 := &RefreshToken{
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * 1),
		UserID:     primitive.NewObjectID(),
		Token:      GetRefreshTokenRepository().FindUnusedToken(),
	}
	GetRefreshTokenRepository().Create(t1)
	t2 := &RefreshToken{
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * -1),
		UserID:     primitive.NewObjectID(),
		Token:      GetRefreshTokenRepository().FindUnusedToken(),
	}
	GetRefreshTokenRepository().Create(t2)

	GetRefreshTokenRepository().CleanUp()

	t1 = GetRefreshTokenRepository().GetOne(t1.ID.Hex())
	t2 = GetRefreshTokenRepository().GetOne(t2.ID.Hex())

	if t1 == nil {
		t.Error("Expected t1 not to be nil")
	}
	if t2 != nil {
		t.Error("Expected t2 to be nil")
	}
}

func TestRefreshTokenGetOne(t *testing.T) {
	clearTestDB()

	t1 := &RefreshToken{
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * 1),
		UserID:     primitive.NewObjectID(),
		Token:      GetRefreshTokenRepository().FindUnusedToken(),
	}
	GetRefreshTokenRepository().Create(t1)

	t1 = GetRefreshTokenRepository().GetOne(t1.ID.Hex())
	if t1 == nil {
		t.Error("Expected t1 not to be nil")
	}
}

func TestRefreshTokenGetOneExpired(t *testing.T) {
	clearTestDB()

	t1 := &RefreshToken{
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * -1),
		UserID:     primitive.NewObjectID(),
		Token:      GetRefreshTokenRepository().FindUnusedToken(),
	}
	GetRefreshTokenRepository().Create(t1)

	t1 = GetRefreshTokenRepository().GetOne(t1.ID.Hex())
	if t1 != nil {
		t.Error("Expected t1 to be nil")
	}
}

func TestRefreshTokenGetByToken(t *testing.T) {
	clearTestDB()

	token := GetRefreshTokenRepository().FindUnusedToken()
	t1 := &RefreshToken{
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * 1),
		UserID:     primitive.NewObjectID(),
		Token:      token,
	}
	GetRefreshTokenRepository().Create(t1)

	t1 = GetRefreshTokenRepository().GetByToken(token)
	if t1 == nil {
		t.Error("Expected t1 not to be nil")
	}
}

func TestRefreshTokenGetByTokenExpired(t *testing.T) {
	clearTestDB()

	token := GetRefreshTokenRepository().FindUnusedToken()
	t1 := &RefreshToken{
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * -1),
		UserID:     primitive.NewObjectID(),
		Token:      token,
	}
	GetRefreshTokenRepository().Create(t1)

	t1 = GetRefreshTokenRepository().GetByToken(token)
	if t1 != nil {
		t.Error("Expected t1 to be nil")
	}
}
