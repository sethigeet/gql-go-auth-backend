init:
	go run github.com/99designs/gqlgen init

generate:
	go generate ./...

run:
	go run .

build:
	go build . -o bin/

.PHONY: init generate
