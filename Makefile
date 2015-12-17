.PHONY: all run build install

all: build

run:
	@go run gomon.go

build:
	@go build gomon.go

install:
	@./install-gomon.sh
