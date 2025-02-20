
test:
	go test ./...

build:
	go build -o bin/go-blog-api .

run:
	go run .

deps:
	go mod tidy
