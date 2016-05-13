#!/bin/bash
set -e

function build-binary {
    echo "--- Building $1/$2"

    ./scripts/utils/build-binary.sh $1 $2
}

echo '--- Setting up $GOPATH'
export GOPATH="$GOPATH:$(pwd)/vendor"
echo $GOPATH

rm -rf pkg

build-binary "windows" "amd64"
build-binary "windows" "386"
build-binary "linux" "amd64"
build-binary "linux" "386"
build-binary "darwin" "amd64"
build-binary "darwin" "386"
build-binary "freebsd" "amd64"
build-binary "freebsd" "386"
