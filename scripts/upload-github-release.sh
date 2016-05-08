#!/bin/bash

VERSION=$(< ./VERSION)
FILES=./releases/*

function github_release() {
    $GOPATH/bin/github-release "$@"
}

echo "--- Creating GitHub release v$VERSION"

github_release release \
    --user "leominov" \
    --repo "self-monitoring" \
    --tag "v$VERSION" \
    --name "Gomon version $VERSION" \
    --description "See README.md" \
    --pre-release

echo "--- Uploading file for release v$VERSION"

for fullfile in $FILES
do
    filename=$(basename "$fullfile")
    echo $filename

    github_release upload \
        --user "leominov" \
        --repo "self-monitoring" \
        --tag "v$VERSION" \
        --name "$filename" \
        --file "$fullfile"

    echo "Done."
done
