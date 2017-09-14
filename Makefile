test:
	@go test

race_test:
	@go test -v -race ./...
