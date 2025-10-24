.PHONY: run
run:
	go run ./cmd/example_one/main.go

.PHONY: example
example:
	http :8000 host=boomatang.com
