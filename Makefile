.PHONY: deps build clean serverless-deploy serverless-remove
deps:
	go mod download
	go mod tidy

build: deps
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bin/searchbypostalcode ./cmd/lambda/searchbypostalcode/main.go

clean:
	rm -rf ./bin

serverless-deploy: clean build
	npm run deploy

serverless-remove: clean build
	npm run remove
