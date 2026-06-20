# Finvera BE - Makefile
# Helper commands untuk development

.PHONY: dev swagger build run tidy

## dev: Jalankan server dengan Air (hot-reload)
dev:
	air

## swagger: Generate ulang swagger docs
swagger:
	swag init -g cmd/server/main.go -o docs

## build: Build binary untuk production
build:
	go build -o bin/finvera-be ./cmd/server/main.go

## run: Jalankan binary
run:
	./bin/finvera-be

## tidy: Bersihkan dan update go.mod
tidy:
	go mod tidy

## help: Tampilkan daftar commands
help:
	@echo "Available commands:"
	@grep -E '^##' Makefile | sed 's/## /  make /'
