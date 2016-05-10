#!/bin/bash
set -e

VERSION=$(< ./VERSION)
PLATFORM=`uname -s`

function build() {
    echo "--- Building release for: $1"

    ./scripts/utils/build-github-release.sh $1
}

function easy_install_script() {
    echo "--- Creating easy install script"
    cp ./scripts/templates/assistant.sh ./releases/assistant.sh
    if [[ $PLATFORM == "Darwin" ]]; then
        sed -i '' "s/TMP_VERSION/$VERSION/g" ./releases/assistant.sh
    else
        sed -i "s/TMP_VERSION/$VERSION/g" ./releases/assistant.sh
    fi
}

export -f build

rm -rf releases tmp

ls pkg/* | xargs -I {} bash -c "build {}"
easy_install_script
