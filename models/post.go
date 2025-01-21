package models

import "time"

type Post struct {
	Id          string    `bson:"_id"`
	PostContent string    `bson:"post_content"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
