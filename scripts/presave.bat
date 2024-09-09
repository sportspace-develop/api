gofmt .\internal\.. .\pkg\.. .\cmd\..
goimports -local "sport-space" -w .\internal\.. .\pkg\.. .\cmd\..
go mod tidy
go test ./...