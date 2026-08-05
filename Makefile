APP_NAME     := devstack
BINARY       := devstack.exe
VERSION      := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE   := $(shell date +%Y-%m-%d)
LDFLAGS      := -ldflags="-s -w -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)"
GO           := go

.PHONY: all build run run-audit run-list test test-short lint install clean help

all: test build

build:
	$(GO) build $(LDFLAGS) -o $(BINARY) .
	@echo "✔ Binário gerado: $(BINARY) (versão: $(VERSION))"

run:
	$(GO) run . bootstrap \
		--stack "Go Backend + React Frontend" \
		--project-name demo-app \
		--dry-run

run-audit:
	$(GO) run . audit --stack go-react

run-list:
	$(GO) run . list-stacks

test:
	$(GO) test ./... -v -race -count=1

test-short:
	$(GO) test ./... -count=1

lint:
	golangci-lint run ./...

install:
	$(GO) install $(LDFLAGS) .
	@echo "✔ devstack instalado em $(shell go env GOPATH)/bin"

clean:
	del /f $(BINARY) 2>nul || true

help:
	@echo "Alvos disponíveis:"
	@grep -E '^## ' Makefile | sed 's/## /  /'
