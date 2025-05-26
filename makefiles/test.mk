.PHONY: test

test:
	# Run all unit tests with coverage
	go test ./... -coverprofile=coverage.out -timeout 30s
	@echo "Coverage report generated at coverage.out"
