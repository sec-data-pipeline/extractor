run: build
	@./bin/extractor

build:
	@go build -o ./bin/extractor
