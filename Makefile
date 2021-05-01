
.PHONY: test
test:
	go test ./... -v $(TESTARGS)  -coverprofile=coverage.out

.PHONY: test-cover
test-cover: test
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: fmt
fmt:
	go fmt ./...
