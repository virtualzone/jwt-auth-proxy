package main

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email          string             `json:"email" bson:"email"`
	HashedPassword string             `json:"password,omitempty" bson:"password"`
	Confirmed      bool               `json:"confirmed" bson:"confirmed"`
	Enabled        bool               `json:"enabled" bson:"enabled"`
	CreateDate     time.Time          `json:"createDate" bson:"createDate"`
	Data           interface{}        `json:"data" bson:"data,omitempty"`
}

type UserRepository struct {
}

var _userRepositoryInstance *UserRepository
var _userRepositoryOnce sync.Once

func GetUserRepository() *UserRepository {
	_userRepositoryOnce.Do(func() {
		_userRepositoryInstance = &UserRepository{}
		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
		mod := mongo.IndexModel{
			Keys: bson.M{
				"email": 1,
			},
			Options: options.Index().SetUnique(true),
		}
		_userRepositoryInstance.GetCollection().Indexes().CreateOne(ctx, mod)
	})
	return _userRepositoryInstance
}

func (r *UserRepository) GetCollection() *mongo.Collection {
	return GetDatatabase().Database.Collection("users")
}

func (r *UserRepository) Create(u *User) {
	res, err := r.GetCollection().InsertOne(context.TODO(), u)
	if err != nil {
		log.Println(err)
	}
	u.ID = res.InsertedID.(primitive.ObjectID)
}

func (r *UserRepository) GetOne(id string) *User {
	var user User
	err := r.GetCollection().FindOne(context.TODO(), GetDatatabase().GetIDFilter(id)).Decode(&user)
	if err != nil {
		return nil
	}
	return &user
}

func (r *UserRepository) GetByEmail(email string) *User {
	var user User
	err := r.GetCollection().FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil
	}
	return &user
}

func (r *UserRepository) Update(u *User) {
	_, err := r.GetCollection().UpdateOne(context.TODO(), bson.M{"_id": u.ID}, bson.M{"$set": u})
	if err != nil {
		log.Println(err)
	}
}

func (r *UserRepository) Delete(u *User) {
	_, err := r.GetCollection().DeleteOne(context.TODO(), bson.M{"_id": u.ID})
	if err != nil {
		log.Println(err)
	}
}

func (r *UserRepository) GetHashedPassword(password string) string {
	pwHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pwHash)
}

func (r *UserRepository) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
