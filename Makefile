BINARY_NAME=build/ttfm
LATEST_TAG="v0.4.0"
COMMON_BUILD_OPTS=-ldflags="-X 'main.Version=${LATEST_TAG}' -s -w"

all: build test

build:
	go build ${COMMON_BUILD_OPTS} -o ${BINARY_NAME}_bot main.go

rpi:
	GOOS=linux GOARCH=arm GOARM=7 go build ${COMMON_BUILD_OPTS} -o ${BINARY_NAME}_rpi main.go

test:
	go test ./...

clean:
	rm build/*
	go clean