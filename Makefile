run: build
	@./bin/filing-extractor

build:
	@go build -o ./bin/filing-extractor
