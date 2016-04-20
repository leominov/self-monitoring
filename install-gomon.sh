#!/bin/bash

GOMON_WD="/opt/self-monitoring/"

GOMON_INIT="/etc/init.d/gomon"
GOMON_INIT_SOURCE="init/gomon"
GOMON_CONFIG="config.json"
GOMON_BINARY="gomon"

export GOMON_WD

log_info() {
    echo "-----> $*"
}

log_verbose() {
    echo "       $*"
}

log_warn() {
    echo " !     $*"
}

log_fail() {
    echo "$@" 1>&2
    exit 1
}

log_info "Gomon Installer"

if [ ! $GOPATH ]; then
    log_fail "ERROR: Missing GOPATH; please see https://golang.org/doc/code.html#GOPATH"
fi

if [[ $EUID -ne 0 ]]; then
    log_fail "ERROR: Must be run with root privileges."
fi

log_info "Building Gomon"
make build
log_verbose "Done"

if [ ! -d "$GOMON_WD" ]; then
    log_verbose "Making directory: $GOMON_WD"
    mkdir -p $GOMON_WD
    log_verbose "Done"
fi

if [ -f "$GOMON_INIT" ]; then
    log_info "Stopping Gomon"
    $GOMON_INIT stop
    log_verbose "Done"
else
    log_info "Coping init script"
    cp $GOMON_INIT_SOURCE $GOMON_INIT
    update-rc.d gomon defaults
    log_verbose "Done"
fi

log_info "Configuring"
if [ ! -f "$GOMON_WD/$GOMON_CONFIG" ]; then
    log_verbose "Coping config file"
    cp $GOMON_CONFIG $GOMON_WD
    log_verbose "Done"
else
    log_warn "Please do not forget to update your monitoring configuration and restart the monitoring:"
    log_warn "$GOMON_INIT restart"
fi

log_info "Coping binary file"
cp gomon $GOMON_WD
log_verbose "Done"

log_info "Starting Gomon"
$GOMON_INIT start
log_verbose "Done"
