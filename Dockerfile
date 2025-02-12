FROM golang:1.23-bookworm

ENV GOLANGCI_LINT_VERSION=v1.64.2
ENV YAEGI_VERSION=v0.16.1
ENV CGO_ENABLED=0

# Install golangci-lint (go linter)
# Cf: https://golangci-lint.run/usage/install
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

# Install Yaegi (go interpreter used by Traefik)
RUN curl -sfL https://raw.githubusercontent.com/traefik/yaegi/master/install.sh | bash -s -- -b $(go env GOPATH)/bin ${YAEGI_VERSION}

WORKDIR /go/src/github.com/ZeroGachis/traefik-block-terminated-clients

COPY ./ ./

RUN \
    go env -w GOPATH=/go && \
    go mod tidy && \
    go mod download && \
    go mod vendor
