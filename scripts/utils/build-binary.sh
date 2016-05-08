#!/bin/bash
set -e

if [[ ${#} -lt 2 ]]
then
    echo "Usage: ${0} [platform] [arch]" >&2
    exit 1
fi

export GOOS=${1}
export GOARCH=${2}

NAME="gomon"

BUILD_PATH="pkg"
BINARY_FILENAME="$NAME-$GOOS-$GOARCH"

echo -e "Building $NAME with:\n"

echo "GOOS=$GOOS"
echo "GOARCH=$GOARCH"
if [[ -n "$GOARM" ]]; then
    echo "GOARM=$GOARM"
fi
echo ""

mkdir -p $BUILD_PATH
go build -v -tags "autogen" -o $BUILD_PATH/$BINARY_FILENAME *.go

chmod +x $BUILD_PATH/$BINARY_FILENAME

echo -e "\nDone: \033[33m$BUILD_PATH/$BINARY_FILENAME\033[0m 💪"
