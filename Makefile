build:
	@go build -o bin/proxied

run: build
	@./bin/proxied

test: 
	@go test ./...  -v

