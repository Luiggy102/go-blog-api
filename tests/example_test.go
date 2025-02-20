package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Luiggy102/go-blog-api/app"
)

func Test_CreatePost(t *testing.T) {

	dbName := "production_blog_" + strconv.Itoa(time.Now().Nanosecond())
	fmt.Println("dbName:", dbName)

	a, err := app.Bootstrap(&app.Config{
		DatabaseUrl: "mongodb://localhost:27017/" + dbName,
	})
	if err != nil {
		t.Error(err)
	}

	s := httptest.NewServer(a.Handler)

	t.Run("Create entry", func(t *testing.T) {
		body, err := json.Marshal(map[string]any{
			"post_content": "Hello World",
		})
		if err != nil {
			t.Error(err)
		}

		req, err := http.NewRequest("POST", s.URL+"/posts", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)

		result := map[string]any{}

		json.NewDecoder(res.Body).Decode(&result)

		if result["post_content"] != "Hello World" {
			t.Errorf("post_content = '%v', want '%v'", result["post_content"], "Hello World")
		}
	})

	t.Run("List entries", func(t *testing.T) {
		req, err := http.NewRequest("GET", s.URL+"/posts", nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.DefaultClient.Do(req)

		result := []any{}

		json.NewDecoder(res.Body).Decode(&result)

		if len(result) != 1 {
			t.Errorf("len = '%v', want '%v'", len(result), 1)
		}
	})

}
