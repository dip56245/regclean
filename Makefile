default: help

## help - show help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## lint - exec golint
lint:
	golangci-lint run

## build - building software
build: lint
	go build -o regclean main.go