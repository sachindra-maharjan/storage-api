.PHONY: build
.PHONY: start

build:
	@echo "Building the datasync service"
	go build -o bin/datasync-service cmd/start/main.go

start:
	@echo "Running the datasync service"
	./bin/datasync-service