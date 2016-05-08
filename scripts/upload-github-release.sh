#!/bin/bash

VERSION=$(< ./VERSION)
FILES=./releases/*

function github_release() {
    $GOPATH/bin/github-release "$@"
}

function notice_release() {
    ./scripts/utils/notice-github-release.sh $1
}

echo "--- Creating GitHub release v$VERSION"

github_release release \
    --user "leominov" \
    --repo "self-monitoring" \
    --tag "v$VERSION" \
    --name "Gomon version $VERSION" \
    --description "See CHANGES.md" \
    --pre-release

if [ $? -eq 0 ]; then
    notice_release $VERSION
    echo "Done."
fi

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

    if [ $? -eq 0 ]; then
        echo "Done."
    fi
done
