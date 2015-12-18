.PHONY: all run build install autogen

all: build

run: autogen
	@go run -tags "autogen" gomon.go

build: autogen
	@go build -tags "autogen" gomon.go

install:
	@./install-gomon.sh

autogen:
	@./.autogen
