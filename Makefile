BINARY=codemind

.PHONY: build docker-build run

build:
	CGO_ENABLED=1 go build -tags treesitter -o $(BINARY) ./

docker-build:
	docker build -t codemind:latest .

run:
	./$(BINARY) $(ARGS)
