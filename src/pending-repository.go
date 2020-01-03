package main

import (
	"context"
	"log"
	"sync"
	"time"

	guuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const PendingActionTypeConfirmAccount = 1
const PendingActionTypeChangeEmail = 2
const PendingActionTypeInitPasswordReset = 3

type PendingAction struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId" bson:"userId"`
	Token      string             `json:"token" bson:"token"`
	ActionType int                `json:"actionType" bson:"actionType"`
	Payload    string             `json:"payload" bson:"payload"`
	CreateDate time.Time          `json:"createDate" bson:"createDate"`
	ExpiryDate time.Time          `json:"expiryDate" bson:"expiryDate"`
}

type PendingActionRepository struct {
}

var _pendingActionRepositoryInstance *PendingActionRepository
var _pendingActionRepositoryOnce sync.Once

func GetPendingActionRepository() *PendingActionRepository {
	_pendingActionRepositoryOnce.Do(func() {
		_pendingActionRepositoryInstance = &PendingActionRepository{}
	})
	return _pendingActionRepositoryInstance
}

func (r *PendingActionRepository) GetCollection() *mongo.Collection {
	return GetDatatabase().Database.Collection("pending_actions")
}

func (r *PendingActionRepository) Create(u *PendingAction) {
	res, err := r.GetCollection().InsertOne(context.TODO(), u)
	if err != nil {
		log.Println(err)
	}
	u.ID = res.InsertedID.(primitive.ObjectID)
}

func (r *PendingActionRepository) GetOne(id string) *PendingAction {
	var pendingAction PendingAction
	err := r.GetCollection().FindOne(context.TODO(), GetDatatabase().GetIDFilter(id)).Decode(&pendingAction)
	if err != nil {
		return nil
	}
	if pendingAction.ExpiryDate.Before(time.Now()) {
		r.Delete(&pendingAction)
		return nil
	}
	return &pendingAction
}

func (r *PendingActionRepository) GetByToken(token string) *PendingAction {
	var pendingAction PendingAction
	err := r.GetCollection().FindOne(context.TODO(), bson.M{"token": token}).Decode(&pendingAction)
	if err != nil {
		return nil
	}
	if pendingAction.ExpiryDate.Before(time.Now()) {
		r.Delete(&pendingAction)
		return nil
	}
	return &pendingAction
}

func (r *PendingActionRepository) GetByPayload(payload string) []*PendingAction {
	var results []*PendingAction
	cur, err := r.GetCollection().Find(context.TODO(), bson.M{
		"payload":    payload,
		"expiryDate": bson.M{"$gte": time.Now()},
	})
	if err != nil {
		return results
	}
	for cur.Next(context.TODO()) {
		var pendingAction PendingAction
		err := cur.Decode(&pendingAction)
		if err != nil {
			return results
		}
		results = append(results, &pendingAction)
	}
	cur.Close(context.TODO())
	return results
}

func (r *PendingActionRepository) Delete(u *PendingAction) {
	_, err := r.GetCollection().DeleteOne(context.TODO(), bson.M{"_id": u.ID})
	if err != nil {
		log.Println(err)
	}
}

func (r *PendingActionRepository) DeleteAllForUser(userID string) {
	_, err := r.GetCollection().DeleteMany(context.TODO(), bson.M{"userId": GetDatatabase().GetObjectID(userID)})
	if err != nil {
		log.Println(err)
	}
}

func (r *PendingActionRepository) FindUnusedToken() string {
	var token string = ""
	for i := 1; i <= 20 && token == ""; i++ {
		token = guuid.New().String()
		if r.GetByToken(token) != nil {
			token = ""
		}
	}
	return token
}

func (r *PendingActionRepository) CleanUp() {
	_, err := r.GetCollection().DeleteMany(context.TODO(), bson.M{"expiryDate": bson.M{"$lte": time.Now()}})
	if err != nil {
		log.Println(err)
	}
}
