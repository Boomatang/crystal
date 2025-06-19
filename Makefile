.PHONY: run
run:
	go run main.go

.PHONY: example
example:
	http :8000 host=boomatang.com
