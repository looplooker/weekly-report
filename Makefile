# Simple Makefile for a Go project

build:
	@echo "Building..."
	
	
	@go build -o weekly-export.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

.PHONY: build run
