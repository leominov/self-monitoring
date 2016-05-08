#!/bin/bash
set -e

function build() {
    echo "--- Building release for: $1"

    VERSION=$(< ./VERSION)

    ./scripts/utils/build-github-release.sh $1 "v$VERSION"
}

export -f build

rm -rf releases

ls pkg/* | xargs -I {} bash -c "build {}"
