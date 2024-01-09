build:
	@go build -o ./bin/zonewatchman -race

run: build
	./bin/zonewatchman