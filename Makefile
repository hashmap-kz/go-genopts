.PHONY: clean test build run

APP_NAME=go-genopts
BUILD_DIR=$(PWD)/bin

# run go only local machine
clean:
	rm -rf ./build

fmt:
	go fmt ./...

test: clean
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: fmt test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) .

run: build
	$(BUILD_DIR)/$(APP_NAME)

run-test:
	go test -v -cover ./...

