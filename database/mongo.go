package database

import (
	"context"
	"log"
	"os"

	"github.com/Luiggy102/go-blog-api/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// insert post
func (mongo *MongoDb) InsertPost(post models.Post) (err error) {
	// the collection
	coll := mongo.db.Database("go_blog").Collection("posts")
	// insert document
	_, err = coll.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}

// get posts
func (mongo *MongoDb) GetPosts(page int64) ([]models.Post, error) {
	coll := mongo.db.Database("go_blog").Collection("posts")

	// find options
	filter := bson.D{{}}
	opts := options.Find().SetLimit(3).SetSkip((page - 1) * 3)

	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Println("error finding documents")
		return nil, err
	}

	results := []models.Post{}

	// send the results
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Println("error parsing results")
		return nil, err
	}
	return results, nil
}

// getPostById
func (mongo *MongoDb) GetPostById(id string) (models.Post, error) {
	var err error
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Post{}, err
	}
	coll := mongo.db.Database("go_blog", nil).Collection("posts", nil)
	filter := bson.D{{Key: "_id", Value: objId}}

	result := coll.FindOne(context.TODO(), filter, nil)
	err = result.Err()
	if err != nil {
		return models.Post{}, err
	}

	post := models.Post{}
	err = result.Decode(&post)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

// UpdatePost
func (mongo *MongoDb) UpdatePost(post models.Post) error {
	coll := mongo.db.Database("go_blog", nil).Collection("posts", nil)
	filter := bson.D{{Key: "_id", Value: post.Id}}
	// log.Println(filter)

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "post_content", Value: post.PostContent},
		// {Key: "created_at", Value: post.CreatedAt},
		{Key: "updated_at", Value: post.UpdatedAt},
	}}}
	opts := options.Update().SetUpsert(true)

	// results, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	// log.Println("upserted id", results.UpsertedID)
	// log.Println("upsert count", results.UpsertedCount)
	// log.Println("matched", results.MatchedCount)
	// log.Println("modified", results.ModifiedCount)

	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)

	if err != nil {
		return err
	}
	return nil
}

// DeletePost
func (mongo *MongoDb) DeletePost(id string) error {
	coll := mongo.db.Database("go_blog").Collection("posts")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objId}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	// log.Println(res.DeletedCount)
	return nil
}
