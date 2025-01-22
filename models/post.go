package models

import "time"

type Post struct {
	// Id          string    `bson:"_id"`
	PostContent string    `json:"post_content" bson:"post_content"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
