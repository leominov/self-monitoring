.PHONY: all run

all: build

run:
	@go run *.go

build:
	@go build *.go
