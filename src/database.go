package main

import (
	"context"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var _databaseInstance *Database
var _databaseOnce sync.Once

func GetDatatabase() *Database {
	_databaseOnce.Do(func() {
		_databaseInstance = &Database{}
	})
	return _databaseInstance
}

func (db *Database) connectMongoDb(url, dbName string) {
	log.Println("Connecting to MongoDB at", url, "...")
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	log.Println("Connected to MongoDB!")
	db.Client = client
	db.Database = client.Database(dbName)
}

func (db *Database) disconnect() {
	log.Println("Closing MongoDB connection...")
	db.Client.Disconnect(context.TODO())
	log.Println("Closed MongoDB connection!")
}

func (db *Database) GetObjectID(id string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(id)
	return objID
}

func (db *Database) GetIDFilter(id string) primitive.M {
	objID := db.GetObjectID(id)
	return bson.M{"_id": objID}
}
