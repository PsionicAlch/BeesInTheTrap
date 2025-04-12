PROJECT_NAME=BeesInTheTrap

build:
	@go build -ldflags="-s -w" -o ./tmp/$(PROJECT_NAME) ./cmd

run: build
	@./tmp/$(PROJECT_NAME)

test:
	@go test ./...

coverage:
	@go test ./... -coverprofile="./tmp/cover.out"
	@go tool cover -html="./tmp/cover.out"

clean:
	@rm -rf ./tmp