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

ASSISTANT_URL=$(curl -s http://tinyurl.com/api-create.php?url=$ASSISTANT_URL)
DESCRIPTION="See CHANGES.md

Install via assistant:
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

echo "--- Creating assistant"
cp ./scripts/templates/assistant.sh ./releases/assistant.sh
if [[ $PLATFORM == "Darwin" ]]; then
    sed -i '' "s/TMP_VERSION/$VERSION/g" ./releases/assistant.sh
else
    sed -i "s/TMP_VERSION/$VERSION/g" ./releases/assistant.sh
fi

echo "--- Uploading file for release v$VERSION"

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
