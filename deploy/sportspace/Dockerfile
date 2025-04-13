FROM golang:1.23 as build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/sportspace/ ./cmd/sportspace/
COPY docs ./docs
COPY internal ./internal
COPY pkg ./pkg

# Build
WORKDIR /app/cmd/sportspace
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.buildVersion=v0.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git show --oneline -s)'" -o /app/cmd/sportspace

FROM ubuntu:latest

WORKDIR /app

COPY --from=build /app/cmd/sportspace/sportspace /sportspace

EXPOSE 8080

# Run
CMD ["/sportspace"]