go-fmt:
	gofumpt -l -w .

go-lint:
	golangci-lint run ./...

tidy:
	go mod tidy

run:
	go run .