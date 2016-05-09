.PHONY: all run install autogen

all: build

run: autogen
	@go run -tags "autogen" gomon.go

install:
	@./install.sh

autogen:
	@./.autogen

.PHONY: build build-binaries build-release build-all

build: autogen
	@go build -tags "autogen" gomon.go

build-binaries: autogen
	@scripts/build-binaries.sh

build-release:
	@scripts/build-github-release.sh

build-all: build-binaries build-release

.PHONY: upload-release upload-pre-release upload-draft-release

upload-release:
	@scripts/upload-github-release.sh --pre-release

upload-pre-release: upload-release

upload-draft-release:
	@scripts/upload-github-release.sh --draft
