.PHONY: all run

all: build

run:
	@go run gomon.go

build:
	@go build gomon.go
