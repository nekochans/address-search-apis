FROM golang:1.16-alpine3.14

LABEL maintainer="https://github.com/nekochans"

WORKDIR /go/app

COPY . .

ARG GOLANGCI_LINT_VERSION=v1.41.1

RUN set -eux && \
  apk update && \
  apk add --no-cache git curl make && \
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION} && \
  go install golang.org/x/tools/cmd/goimports@latest

ENV CGO_ENABLED 0
