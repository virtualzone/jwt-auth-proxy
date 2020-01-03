package main

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	guuid "github.com/google/uuid"
)

type RefreshToken struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId" bson:"userId"`
	Token      string             `json:"token" bson:"token"`
	CreateDate time.Time          `json:"createDate" bson:"createDate"`
	ExpiryDate time.Time          `json:"expiryDate" bson:"expiryDate"`
}

type RefreshTokenRepository struct {
}

var _refreshTokenRepositoryInstance *RefreshTokenRepository
var _refreshTokenRepositoryOnce sync.Once

func GetRefreshTokenRepository() *RefreshTokenRepository {
	_refreshTokenRepositoryOnce.Do(func() {
		_refreshTokenRepositoryInstance = &RefreshTokenRepository{}
	})
	return _refreshTokenRepositoryInstance
}

func (r *RefreshTokenRepository) GetCollection() *mongo.Collection {
	return GetDatatabase().Database.Collection("refresh_tokens")
}

func (r *RefreshTokenRepository) Create(u *RefreshToken) {
	res, err := r.GetCollection().InsertOne(context.TODO(), u)
	if err != nil {
		log.Println(err)
	}
	u.ID = res.InsertedID.(primitive.ObjectID)
}

func (r *RefreshTokenRepository) GetOne(id string) *RefreshToken {
	var refreshToken RefreshToken
	err := r.GetCollection().FindOne(context.TODO(), GetDatatabase().GetIDFilter(id)).Decode(&refreshToken)
	if err != nil {
		return nil
	}
	if refreshToken.ExpiryDate.Before(time.Now()) {
		r.Delete(&refreshToken)
		return nil
	}
	return &refreshToken
}

func (r *RefreshTokenRepository) GetByToken(token string) *RefreshToken {
	var refreshToken RefreshToken
	err := r.GetCollection().FindOne(context.TODO(), bson.M{"token": token}).Decode(&refreshToken)
	if err != nil {
		return nil
	}
	if refreshToken.ExpiryDate.Before(time.Now()) {
		r.Delete(&refreshToken)
		return nil
	}
	return &refreshToken
}

func (r *RefreshTokenRepository) DeleteAllForUser(userID string) {
	_, err := r.GetCollection().DeleteMany(context.TODO(), bson.M{"userId": GetDatatabase().GetObjectID(userID)})
	if err != nil {
		log.Println(err)
	}
}

func (r *RefreshTokenRepository) Delete(u *RefreshToken) {
	_, err := r.GetCollection().DeleteOne(context.TODO(), bson.M{"_id": u.ID})
	if err != nil {
		log.Println(err)
	}
}

func (r *RefreshTokenRepository) FindUnusedToken() string {
	var token string = ""
	for i := 1; i <= 20 && token == ""; i++ {
		token = guuid.New().String()
		if r.GetByToken(token) != nil {
			token = ""
		}
	}
	return token
}

func (r *RefreshTokenRepository) CleanUp() {
	_, err := r.GetCollection().DeleteMany(context.TODO(), bson.M{"expiryDate": bson.M{"$lte": time.Now()}})
	if err != nil {
		log.Println(err)
	}
}
