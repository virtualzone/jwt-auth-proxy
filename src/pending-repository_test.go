package main

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPendingActionCleanUp(t *testing.T) {
	clearTestDB()

	pa1 := &PendingAction{
		ActionType: 1,
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * 1),
		UserID:     primitive.NewObjectID(),
		Payload:    "",
		Token:      GetPendingActionRepository().FindUnusedToken(),
	}
	GetPendingActionRepository().Create(pa1)
	pa2 := &PendingAction{
		ActionType: 1,
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * -1),
		UserID:     primitive.NewObjectID(),
		Payload:    "",
		Token:      GetPendingActionRepository().FindUnusedToken(),
	}
	GetPendingActionRepository().Create(pa2)

	GetPendingActionRepository().CleanUp()

	pa1 = GetPendingActionRepository().GetOne(pa1.ID.Hex())
	pa2 = GetPendingActionRepository().GetOne(pa2.ID.Hex())

	if pa1 == nil {
		t.Error("Expected pa1 not to be nil")
	}
	if pa2 != nil {
		t.Error("Expected pa1 to be nil")
	}
}

func TestPendingActionGetOne(t *testing.T) {
	clearTestDB()

	pa1 := &PendingAction{
		ActionType: 1,
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * 1),
		UserID:     primitive.NewObjectID(),
		Payload:    "",
		Token:      GetPendingActionRepository().FindUnusedToken(),
	}
	GetPendingActionRepository().Create(pa1)

	pa1 = GetPendingActionRepository().GetOne(pa1.ID.Hex())
	if pa1 == nil {
		t.Error("Expected pa1 not to be nil")
	}
}

func TestPendingActionGetOneExpired(t *testing.T) {
	clearTestDB()

	pa1 := &PendingAction{
		ActionType: 1,
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * -1),
		UserID:     primitive.NewObjectID(),
		Payload:    "",
		Token:      GetPendingActionRepository().FindUnusedToken(),
	}
	GetPendingActionRepository().Create(pa1)

	pa1 = GetPendingActionRepository().GetOne(pa1.ID.Hex())
	if pa1 != nil {
		t.Error("Expected pa1 to be nil")
	}
}

func TestPendingActionGetByToken(t *testing.T) {
	clearTestDB()

	token := GetPendingActionRepository().FindUnusedToken()
	pa1 := &PendingAction{
		ActionType: 1,
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * 1),
		UserID:     primitive.NewObjectID(),
		Payload:    "",
		Token:      token,
	}
	GetPendingActionRepository().Create(pa1)

	pa1 = GetPendingActionRepository().GetByToken(token)
	if pa1 == nil {
		t.Error("Expected pa1 not to be nil")
	}
}

func TestPendingActionGetByTokenExpired(t *testing.T) {
	clearTestDB()

	token := GetPendingActionRepository().FindUnusedToken()
	pa1 := &PendingAction{
		ActionType: 1,
		CreateDate: time.Now(),
		ExpiryDate: time.Now().Add(time.Duration(time.Minute) * -1),
		UserID:     primitive.NewObjectID(),
		Payload:    "",
		Token:      token,
	}
	GetPendingActionRepository().Create(pa1)

	pa1 = GetPendingActionRepository().GetByToken(token)
	if pa1 != nil {
		t.Error("Expected pa1 to be nil")
	}
}
