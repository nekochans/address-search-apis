.PHONY: deps build clean serverless-deploy serverless-remove lint format test test-ci ci
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

lint:
	go vet ./...
	golangci-lint run ./...

format:
	gofmt -l -s -w .
	goimports -w -l ./

test:
	go clean -testcache
	go test -p 1 -v $$(go list ./... | grep -v /node_modules/)

test-ci:
	go clean -testcache
	go test -p 1 -v -coverprofile coverage.out -covermode atomic $$(go list ./... | grep -v /node_modules/)

ci: lint test-ci
	go mod tidy && git diff -s --exit-code go.sum
