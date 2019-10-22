all: build

build: fmt
	go build
	zip git-webhook-lambda.zip git-webhook-lambda

fmt:
	go fmt $(shell find . -type d | egrep -v "(terraform|vendor|.git)")

include Makefile.vars

run: fmt build
	debug=true rules=$(RULES) ./git-webhook-lambda -local
