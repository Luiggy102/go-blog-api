package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/Luiggy102/go-blog-api/models"
)

type MongoDb struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDb(db_url string) (*MongoDb, error) {

	cs, err := connstring.Parse(db_url)
	if err != nil {
		return nil, err
	}
	err = cs.Validate()
	if err != nil {
		return nil, err
	}

	optionsClient := options.Client().ApplyURI(db_url)
	client, err := mongo.Connect(context.Background(), optionsClient)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	log.Println("Database connected successfully!")

	return &MongoDb{
		client: client,
		db:     client.Database(cs.Database),
	}, nil
}

func (mongo *MongoDb) Close() error {
	err := mongo.client.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

// insert post
func (mongo *MongoDb) InsertPost(post models.Post) (err error) {
	// the collection
	coll := mongo.db.Collection("posts")
	// insert document
	_, err = coll.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}

// get posts
func (mongo *MongoDb) GetPosts(page int64) ([]models.Post, error) {
	coll := mongo.db.Collection("posts")

	// find options
	filter := bson.D{{}}

	var opts *options.FindOptions
	page_size := int64(10)

	// 0 = all posts, 1... normal pagination
	if page == 0 {
		opts = options.Find()
	} else {
		opts = options.Find().SetLimit(page_size).SetSkip((page - 1) * page_size)
	}

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
	coll := mongo.db.Collection("posts", nil)
	filter := bson.D{{Key: "_id", Value: id}}

	result := coll.FindOne(context.TODO(), filter, nil)
	err := result.Err()
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
	coll := mongo.db.Collection("posts", nil)
	filter := bson.D{{Key: "_id", Value: post.Id}}
	// log.Println(filter)

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "post_title", Value: post.PostTitle},
				{Key: "post_content", Value: post.PostContent},
				{Key: "updated_at", Value: post.UpdatedAt},
			},
		},
		{
			Key: "$setOnInsert",
			Value: bson.D{
				{Key: "created_at", Value: post.CreatedAt},
			},
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)

	if err != nil {
		return err
	}
	return nil
}

// DeletePost
func (mongo *MongoDb) DeletePost(id string) error {
	coll := mongo.db.Collection("posts")
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	// log.Println(res.DeletedCount)
	return nil
}
