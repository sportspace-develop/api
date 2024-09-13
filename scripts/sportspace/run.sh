swag init -o ./docs -g ./internal/adapter/api/rest/rest.go
swag fmt
go run -ldflags "-X main.buildVersion=v0.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git show --oneline -s)'" ./cmd/sportspace/sportspace.go
