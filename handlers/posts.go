package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/Luiggy102/go-blog-api/database"
	"github.com/Luiggy102/go-blog-api/models"
)

type upsertPostRequest struct {
	PostContent string `json:"post_content"`
}
type insertPostResponse struct {
	// Id interface{} `json:"post_id"`
	// Message string      `json:"message"`
	models.Post
}
type messageResponse struct {
	Message string `json:"message"`
}

func InsertPostHandler(mongo *database.MongoDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// decode request
		postRequest := upsertPostRequest{}
		err := json.NewDecoder(r.Body).Decode(&postRequest)
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		if err != nil {
			log.Println("Error decoding body:", err.Error())
			http.Error(w, "JSON is not valid", http.StatusInternalServerError)
			return
		}

		// create a random id
		newId := uuid.New().String()

		// new post
		now := time.Now()
		p := models.Post{
			Id:          newId,
			PostContent: postRequest.PostContent,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// insert post into db
		err = mongo.InsertPost(p)
		if err != nil {
			log.Println("Error inserting post:", err.Error())
			http.Error(w, "Error persisting data", http.StatusInternalServerError)
			return
		}
		log.Printf("Post created with ID: %v", newId)

		// send reponse
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&insertPostResponse{
			Post: p,
		})
		if err != nil {
			log.Println("DEV: Error encoding response:", err.Error())
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
	}
}

func GetPostsHandler(mongo *database.MongoDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		// pagination
		pageStr := r.URL.Query().Get("page")
		// default value
		var page = uint64(1)

		if pageStr != "" {
			// not empty = now value
			page, err = strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				// an error like `/posts?page=foo`
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		posts, err := mongo.GetPosts(int64(page))
		if err != nil {
			log.Fatalln("Error fetching data from db")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// send the response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Fatalln("error parsing posts data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetPostsbyIdHandler
func GetPostsbyIdHandler(mongo *database.MongoDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get id
		id := r.PathValue("id")
		post, err := mongo.GetPostById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// return post
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// UpdatePost
func UpdatePostHander(mongo *database.MongoDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle Request

		// decode request
		updateRequest := upsertPostRequest{}
		err := json.NewDecoder(r.Body).Decode(&updateRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get post Id
		id := r.PathValue("id")

		// update data
		err = mongo.UpdatePost(models.Post{
			Id:          id,
			PostContent: updateRequest.PostContent,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// send response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&messageResponse{"ok"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// delete post
func DeletePostHandler(mongo *database.MongoDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// post id
		postId := r.PathValue("id")

		err := mongo.DeletePost(postId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// send response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&messageResponse{"ok"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
