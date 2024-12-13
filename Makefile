# Simple Makefile for a Go project

build:
	@echo "Building..."
	@go build -o weekly-report.exe cmd/main.go

# Run the application
run:
	@go run cmd/main.go

.PHONY: build run
