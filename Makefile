go-fmt:
	gofumpt -l -w .

go-lint:
	golangci-lint run ./...

tidy:
	go mod tidy

run:
	go run .

api-up:
	docker-compose -f ./deployments/compose/api-server.docker-compose.yaml --env-file=./.env up -d --build

api-down:
	docker-compose -f ./deployments/compose/api-server.docker-compose.yaml --env-file=./.env up down

logging-up:
	docker-compose -f ./deployments/compose/logging.docker-compose.yaml --env-file=./.env up -d --build

logging-down:
	docker-compose -f ./deployments/compose/logging.docker-compose.yaml --env-file=./.env up down

metrics-up:
	docker-compose -f ./deployments/compose/metrics.docker-compose.yaml --env-file=./.env up -d --build

metrics-down:
	docker-compose -f ./deployments/compose/metrics.docker-compose.yaml --env-file=./.env up down

telemetry-up:
	docker-compose -f ./deployments/compose/telemetry.docker-compose.yaml --env-file=./.env up -d --build

telemetry-down:
	docker-compose -f ./deployments/compose/telemetry.docker-compose.yaml --env-file=./.env up down