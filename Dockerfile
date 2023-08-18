FROM golang:1.20 AS build

WORKDIR /usr/local/go/src/sport-space-api

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

WORKDIR /usr/local/go/src/sport-space-api/cmd/server

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app


FROM ubuntu:latest

WORKDIR /app

COPY .env.develop /.env
COPY docs ./docs
COPY --from=build /build/app ./

EXPOSE 8080
ENTRYPOINT ["./app"]