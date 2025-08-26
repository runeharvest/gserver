run: build
	@cd bin && ./gserver

build:
	@mkdir -p bin
	@go build -o bin/gserver ./cmd

init: install-buf install-wire

.PHONY: proto
proto:
	buf generate shared/proto

install-buf:
	@echo "Installing buf..."
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@echo "Verifying buf installation..."
	@buf --version

install-wire:
	@echo "Installing wire..."
	@go install github.com/google/wire/cmd/wire@latest
	@echo "Verifying wire installation..."

.PHONY: wire
wire:
	wire gen ./...