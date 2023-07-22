sec:
	gosec ./..

lint:
	golangci-lint run ./...

test:
	cd ./cmd/benchmark && go run main.go

dep:
	go mod tidy

up:
	cd docker && docker-compose up -d --build

run-api:
	cd cmd/api && go run main.go

run-benchmark:
	cd cmd/benchmark && go run main.go