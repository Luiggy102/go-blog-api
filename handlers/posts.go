package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Luiggy102/go-blog-api/database"
	"github.com/Luiggy102/go-blog-api/models"
)

type insertPostRequest struct {
	PostContent string `json:"post_content"`
}
type insertPostResponse struct {
	Id interface{} `json:"post_id"`
	// Message string      `json:"message"`
	models.Post
}

func InsertPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mongo, err := database.NewMongoDb()
		defer func() {
			err = mongo.Close()
			if err != nil {
				log.Fatalln("error closing db", err.Error())
			}
		}()

		if err != nil {
			log.Fatalln("Error in db connection")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// decode request
		postRequest := insertPostRequest{}
		err = json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// new post
		p := models.Post{
			PostContent: postRequest.PostContent,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// insert post into db
		postId, err := mongo.InsertPost(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Post created with ID: %v", postId)

		// send reponse
		w.Header().Set("Content-Type", "application/json")
		// err = json.NewEncoder(w).Encode(&insertPostResponse{
		// 	Message: "ok",
		// })
		err = json.NewEncoder(w).Encode(&insertPostResponse{
			Id:   postId,
			Post: p,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
