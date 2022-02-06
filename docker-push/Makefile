mod:
	go mod tidy

build: mod
	CGO_ENABLED=0 go build .

run:
	go run main.go

run-package: build
	go run docker-push
