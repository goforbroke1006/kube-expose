NAME=kube-expose

.PHONY: build
build:
	env GOOS=linux GOARCH=amd64 go build -o ./build/${NAME} ./cmd/${NAME}/main.go

install: build
	sudo cp ./build/${NAME} /usr/local/bin/
	sudo chmod +x /usr/local/bin/${NAME}