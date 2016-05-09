#!/bin/bash
# curl -sSL https://github.com/leominov/self-monitoring/releases/download/0.1.1/download.sh | sh

set -e

VERSION="TMP_VERSION"

PLATFORM=`uname -s`
ARCH=`uname -m`

NAME="gomon-$PLATFORM-$ARCH.tar.gz"
URL="https://github.com/leominov/self-monitoring/releases/download/$VERSION/$NAME"

echo "--- Downloading $PLATFORM/$ARCH"
curl -OL $URL
tar -xf $NAME
rm $NAME

echo "--- Preparing for installation $VERSION"
if [ -f "config.json" ]; then
    echo "    Creating new configuration (config.json)"
    cp example.config.json config.json
fi
echo " !  Please do not forget to update your monitoring configuration and restart the monitoring:"
echo "    service gomon restart"

echo "--- Installing"
sudo ./install.sh
