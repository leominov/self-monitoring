#!/bin/bash
set -e

if [[ ${#} -lt 1 ]]
then
    echo "Usage: ${0} [file]" >&2
    exit 1
fi

function info {
    echo -e "\033[35m$1\033[0m"
}

BINARY_PATH=${1}

BASE_DIRECTORY=`pwd`
INIT_DIDECTORY=$BASE_DIRECTORY/init
TMP_DIRECTORY=$BASE_DIRECTORY/tmp
RELEASE_DIRECTORY=$BASE_DIRECTORY/releases

RELEASE_NAME=$(basename "$BINARY_PATH")
RELEASE_NAME="${RELEASE_NAME%.*}"

TMP_RELEASE_DIRECTORY=$TMP_DIRECTORY/$RELEASE_NAME

rm -rf $TMP_RELEASE_DIRECTORY
mkdir -p $TMP_RELEASE_DIRECTORY

if [[ "$BINARY_PATH" == *"Windows"* ]]; then
    RELEASE_FILE_NAME="$RELEASE_NAME.zip"

    info "Copying binary"
    cp $BINARY_PATH $TMP_RELEASE_DIRECTORY/gomon.exe

    info "Copying README.md"
    cp $BASE_DIRECTORY/README.md $TMP_RELEASE_DIRECTORY/

    info "Copying config"
    cp $BASE_DIRECTORY/example.config.json $TMP_RELEASE_DIRECTORY/

    info "Zipping up the files"
    cd $TMP_RELEASE_DIRECTORY
    zip -X -r "../$RELEASE_FILE_NAME" *
else
    RELEASE_FILE_NAME="$RELEASE_NAME.tar.gz"

    info "Copying binary"
    cp $BINARY_PATH $TMP_RELEASE_DIRECTORY/gomon
    chmod +x $TMP_RELEASE_DIRECTORY/gomon

    info "Copying README.md"
    cp $BASE_DIRECTORY/README.md $TMP_RELEASE_DIRECTORY/

    info "Copying config"
    cp $BASE_DIRECTORY/example.config.json $TMP_RELEASE_DIRECTORY/

    info "Copying init script"
    mkdir -p $TMP_RELEASE_DIRECTORY/init/
    cp $INIT_DIDECTORY/gomon $TMP_RELEASE_DIRECTORY/init/
    cp $INIT_DIDECTORY/README.md $TMP_RELEASE_DIRECTORY/init/

    info "Copying install script"
    cp $BASE_DIRECTORY/install.sh $TMP_RELEASE_DIRECTORY/

    info "Tarring up the files"
    cd $TMP_RELEASE_DIRECTORY
    tar cfvz ../$RELEASE_FILE_NAME .
fi

mkdir -p $RELEASE_DIRECTORY
mv $TMP_DIRECTORY/$RELEASE_FILE_NAME $RELEASE_DIRECTORY/

echo -e "👏 Created release \033[33m$RELEASE_DIRECTORY/$RELEASE_FILE_NAME\033[0m"
