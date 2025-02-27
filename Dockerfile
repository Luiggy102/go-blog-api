FROM golang:1.22-alpine AS contructor

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . ./
RUN go build -o /blog-api

FROM alpine:latest

WORKDIR /app
COPY --from=contructor /app/config.json /app/config.json
COPY --from=contructor /blog-api /app/blog-api
