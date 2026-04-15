test:
	go clean -testcache
	go test -v -coverprofile=coverage.out "$(shell go list ./... | grep -v "examples")"
	go tool cover -html=coverage.out -o coverage.html
