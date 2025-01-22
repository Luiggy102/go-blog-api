package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	db *mongo.Client
}

func getDbURl() string {
	var err error = godotenv.Load("./.env")
	if err != nil {
		log.Fatalln("Error loading env vars", err.Error())
	}
	return os.Getenv("DATABASE_URL")
}

func NewMongoDb() (*MongoDb, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(getDbURl()))
	if err != nil {
		return nil, err
	}
	return &MongoDb{db: client}, nil
}

func (mongo *MongoDb) Close() error {
	err := mongo.db.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	return nil
}
