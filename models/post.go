package models

import (
	"time"
)

type Post struct {
	Id          string    `json:"id" bson:"_id"`
	PostTitle   string    `json:"post_title" bson:"post_title"`
	PostContent string    `json:"post_content" bson:"post_content"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
