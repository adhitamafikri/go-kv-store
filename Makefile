DEFAULT_GOAL := help

help:
	@echo "Go KV Store"
	@echo "Usage: make <target>"
	@echo "TBD"

run:
	@go run cmd/kvserver/main.go
