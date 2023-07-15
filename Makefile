sec:
	gosec ./..

lint:
	golangci-lint run ./...

test:
	go test -v ./...