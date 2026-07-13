.PHONY: test dev

build:
	go build -o quitocoin cmd/main.go

run:
	./quitocoin

dev:
	go run cmd/main.go

test:
	go test ./test/*/* -v
