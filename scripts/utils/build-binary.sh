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

echo -e "Building $NAME with:\n"

echo "GOOS=$GOOS"
echo "GOARCH=$GOARCH"
if [[ -n "$GOARM" ]]; then
    echo "GOARM=$GOARM"
fi

pkgOS=$GOOS
case "$pkgOS" in
    darwin)
        pkgOS=Darwin
        ;;
    freebsd)
        pkgOS=FreeBSD
        ;;
    linux)
        pkgOS=Linux
        ;;
    windows)
        pkgOS=Windows
        ;;
    *)
        echo >&2 "Error: can't convert $pkgOS to an appropriate value for 'uname -s'"
        exit 1
        ;;
esac

pkgArch=$GOARCH
case "$pkgArch" in
    amd64)
        pkgArch=x86_64
        ;;
    386)
        pkgArch=i386
        ;;
    arm)
        pkgArch=armel
        ;;
    *)
        echo >&2 "Error: can't convert $pkgArch to an appropriate value for 'uname -m'"
        exit 1
        ;;
esac

echo ""

BINARY_FILENAME="$NAME-$pkgOS-$pkgArch"

if [[ "$GOOS" == "windows" ]]; then
    BINARY_FILENAME="$BINARY_FILENAME.exe"
fi

mkdir -p $BUILD_PATH
go build -v -tags "autogen" -o $BUILD_PATH/$BINARY_FILENAME *.go

chmod +x $BUILD_PATH/$BINARY_FILENAME

echo -e "\nDone: \033[33m$BUILD_PATH/$BINARY_FILENAME\033[0m 💪"
