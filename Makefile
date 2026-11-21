DEFAULT_GOAL := help

help:
	@echo "Go KV Store"
	@echo "Usage: make <target>"
	@echo "TBD"

run:
	@go run cmd/kvserver/main.go

test:
	@echo "Running all tests..."
	@go test ./...

test_stage_01:
	@echo "Running all tests on Stage 01..."
	@go test -run Stage01 -v ./...
