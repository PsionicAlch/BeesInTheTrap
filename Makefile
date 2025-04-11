PROJECT_NAME=BeesInTheTrap

build:
	@go build -ldflags="-s -w" -o ./tmp/$(PROJECT_NAME) ./cmd

run: build
	@./tmp/$(PROJECT_NAME)

clean:
	@rm -rf ./tmp