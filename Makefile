BINARY_NAME=seactl
build:
	go build -o seactl .

run:
	go run main.go


compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-x86 .
	GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME}-aarch64 .
	#GOOS=freebsd GOARCH=386 go build -o ${BINARY_NAME} .

test:
	go test -v ./... -cover

all: test build compile run
