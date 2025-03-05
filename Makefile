go-fmt:
	gofumpt -l -w .

go-lint:
	golangci-lint run ./...

tidy:
	go mod tidy

run:
	go run .

up:
	docker-compose -f ./docker-compose.yaml --env-file=./.env up -d --build

down:
	docker-compose -f ./docker-compose.yaml --env-file=./.env up down