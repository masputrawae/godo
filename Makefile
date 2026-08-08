generate:
	@templ generate

build:
	@go build -o ./bin/app

run: generate build
	@./bin/app

test:
	@go test -v -cover ./...

clean:
	@rm -rf ./bin
