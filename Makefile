BINARY_NAME=micro-cli
VERSION="-X main.VERSION=`cat 'version'` -X main.GIT_HASH=`git rev-parse HEAD` -s"

default:
	@echo 'Usage of make: [ build | linux_build | windows_build | run | clean ]'

build: 
	@go build -ldflags ${VERSION} -o ./bin/${BINARY_NAME} ./

linux_build: 
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags ${VERSION} -o ./bin/${BINARY_NAME} ./

windows_build: 
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${VERSION} -o ./bin/${BINARY_NAME}.exe ./

run: build
	./bin/${BINARY_NAME}

clean: 
	@rm -f ./bin/${BINARY_NAME}*

.PHONY: default build linux_build windows_build run clean