.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/gateway gateway/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
