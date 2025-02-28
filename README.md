# go-blog-api

go-blog-api is a **REST API** project for create, read, update and delete posts (CRUD operations). It uses a **Go** backend and a **Mongo** database for storing the data.
Also can be run with Docker.

## Base url
```
http://localhost:8080/
```

##### Run it with [Go] (Need running MongoDB with a db called *go_blog*).
```
git clone https://github.com/Luiggy102/go-blog-api.git
cd go-blog-api
go run main.go
```

##### Run it with [Docker]
```
git clone https://github.com/Luiggy102/go-blog-api.git
cd go-blog-api
docker compose up
```

## Endpoints
`GET /`: Shows a welcome message

`GET /posts`: List the firsts posts

`GET /posts/{id}`: Get the post info by the post ID

`GET /posts?page=2`: Pagination feature to see posts

`POST /posts`: Insert a post in JSON format
Example:
```json
{
    "post_title": "post title"
    "post_content": "sample text"
}

```
`PUT /posts/{id}`: Update the post using the post ID, receive a JSON with the content to be updated
```json
{
    "post_title": "new post title"
    "post_content": "update text"
}
```

`DELETE /posts/{id}`: Delete the post using the post ID
