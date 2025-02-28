package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Luiggy102/go-blog-api/database"
)

func HomeHandler(mongo *database.MongoDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		posts, err := mongo.GetPosts(0)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error fetching data from db", http.StatusInternalServerError)
			return
		}

		t, err := template.New("home.html").Parse(`

<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Blog</title>
<style>
body {
  max-width: 720px;
  margin: 0 auto;
}
article {
  border: solid silver 1px;
  border-radius: 3px;
  padding: 16px;
  margin: 16px;
}
</style>
</head>
<body>
	<h1>Blog</h1>
	{{ range $post := .posts }}
	<article class="post">
        <h3>{{ .PostTitle}}</h2>
		{{ .PostContent }}
	</article>
	{{ end }}
	
</body>
</html>
`)
		if err != nil {
			log.Println("DEV: home template is wrong:", err.Error())
			http.Error(w, "Oops", http.StatusInternalServerError)
		}

		err = t.Execute(w, map[string]interface{}{
			"posts": posts,
		})
		if err != nil {
			log.Println("DEV: home template is wrong:", err.Error())
		}
	}
}
