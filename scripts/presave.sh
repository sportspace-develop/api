gofmt ./internal/.. ./pkg/.. ./cmd/..
goimports -local "sport-space" -w ./internal/.. ./pkg/.. ./cmd/..
go mod tidy
go test ./...
swag init -o ./docs -g ./internal/adapter/api/rest/rest.go
swag fmt