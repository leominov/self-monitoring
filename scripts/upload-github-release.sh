#!/bin/bash

VERSION=$(< ./VERSION)
FILES=./releases/*
RELEASE_TYPE=$1
PLATFORM=`uname -s`

ASSISTANT_URL="https://github.com/leominov/self-monitoring/releases/download/$VERSION/assistant.sh"

function github_release() {
    $GOPATH/bin/github-release "$@"
}

function notice_release() {
    ./scripts/utils/notice-github-release.sh $1
}

echo "--- Creating GitHub release v$VERSION"

ASSISTANT_TINY_URL=$(curl -s http://tinyurl.com/api-create.php?url=$ASSISTANT_URL)
DESCRIPTION="See CHANGES.md

Quick and easy install via:
curl -sSL $ASSISTANT_TINY_URL | sh
Blocked?

Try:
curl -sSL $ASSISTANT_URL | sh
"

github_release release \
    --user "leominov" \
    --repo "self-monitoring" \
    --tag "$VERSION" \
    --name "$VERSION" \
    --description "$DESCRIPTION" \
    $RELEASE_TYPE

if [ $? -eq 0 ]; then
    notice_release $VERSION
    echo "Done."
fi

echo "--- Uploading files for release v$VERSION"

for fullfile in $FILES
do
    filename=$(basename "$fullfile")
    echo $filename

    github_release upload \
        --user "leominov" \
        --repo "self-monitoring" \
        --tag "$VERSION" \
        --name "$filename" \
        --file "$fullfile"

    if [ $? -eq 0 ]; then
        echo "Done."
    fi
done
