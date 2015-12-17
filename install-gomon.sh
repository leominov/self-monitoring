#!/bin/bash

GOMON_WD="/opt/self-monitoring/"

GOMON_INIT="/etc/init.d/gomon"
GOMON_INIT_SOURCE="init/gomon"
GOMON_CONFIG="config.json"
GOMON_BINARY="gomon"

export GOMON_WD

if [[ $EUID -ne 0 ]]; then
    echo "ERROR: Must be run with root privileges."
    exit 1
fi

if [ ! -f $GOMON_BINARY ]; then
    make build
fi

if [ ! -d $GOMON_WD ]; then
    mkdir -p $GOMON_WD
fi

if [ -f $GOMON_INIT ]; then
    $GOMON_INIT stop
else
    cp $GOMON_INIT_SOURCE $GOMON_INIT
fi

cp gomon $GOMON_WD
cp $GOMON_CONFIG $GOMON_WD

update-rc.d gomon defaults

$GOMON_INIT start
