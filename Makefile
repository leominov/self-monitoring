.PHONY: all run build build-binaries build-release build-all upload-release install autogen

all: build

run: autogen
	@go run -tags "autogen" gomon.go

build: autogen
	@go build -tags "autogen" gomon.go

build-binaries:
	@scripts/build-binaries.sh

build-release:
	@scripts/build-github-release.sh

build-all: build-binaries build-release

upload-release:
	@scripts/upload-github-release.sh

install:
	@./install-gomon.sh

autogen:
	@./.autogen
