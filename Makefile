.PHONY: all run build install autogen

all: build

run:
	@go run gomon.go

build: autogen
	@go build -tags "autogen" gomon.go

install:
	@./install-gomon.sh

autogen:
	@./.autogen
