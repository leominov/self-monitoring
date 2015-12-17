#!/bin/bash

VERSION=$(< ./VERSION)
GITCOMMIT=$(git rev-parse --short HEAD)
BUILDTIME=$(date -u)

cat > gomonversion/version_autogen.go <<DVEOF
// +build autogen

package gomonversion

// This file is auto-generated at build-time
const (
	GitCommit string = "$GITCOMMIT"
	Version   string = "$VERSION"
	BuildTime string = "$BUILDTIME"
)
DVEOF
