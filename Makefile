BINARY_NAME=ttfm
 
all: build test

build:
	go build -o ${BINARY_NAME}_bot main.go

rpi:
	GOOS=linux GOARCH=arm GOARM=7 go build -o ${BINARY_NAME}_rpi main.go

test:
	go test ./...

clean:
	go clean