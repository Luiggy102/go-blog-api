package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Luiggy102/go-blog-api/database"
	"github.com/Luiggy102/go-blog-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		postRequest := upsertPostRequest{}
		err = json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// create a random id
		newId := primitive.NewObjectID()
		// newId, err := ksuid.NewRandom()
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// new post
		p := models.Post{
			Id:          newId,
			PostContent: postRequest.PostContent,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// insert post into db
		err = mongo.InsertPost(p)
		log.Printf("Post created with ID: %v", newId.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// send reponse
		w.Header().Set("Content-Type", "application/json")
		// err = json.NewEncoder(w).Encode(&insertPostResponse{
		// 	Message: "ok",
		// })
		err = json.NewEncoder(w).Encode(&insertPostResponse{
			Post: p,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetPostsHandler() http.HandlerFunc {
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
func GetPostsbyIdHandler() http.HandlerFunc {
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
func UpdatePostHander() http.HandlerFunc {
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

		// handle Request

		// decode request
		updateRequest := upsertPostRequest{}
		err = json.NewDecoder(r.Body).Decode(&updateRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get post Id
		id := r.PathValue("id")
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// update data
		err = mongo.UpdatePost(models.Post{
			Id:          objId,
			PostContent: updateRequest.PostContent,
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
