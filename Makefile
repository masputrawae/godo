generate:
	@templ generate

build:
	@go build -o ./bin/app

test:
	@go test -v -cover ./...

clean:
	@rm -rf ./bin

run: generate build
	@./bin/app
