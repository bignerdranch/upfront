test:
	go test ./...
.PHONY: test

# Runs the example server
example:
	go run ./example/
.PHONY: example
