#!/usr/bin/env bash
PKG_NAME="podman-tui"
GIT_TOPDIR=$(git rev-parse --show-toplevel)
GIT_LAST_COMMIT=$(git log -1 --format=%H)
SPEC_FILE=${GIT_TOPDIR}/rpm/${PKG_NAME}.spec

VERSION=$(grep -E "[[:space:]]+appVersion\s*=" ${GIT_TOPDIR}/cmd/version.go  | cut -d\" -f2 | sed 's/-//')

sed -i "s/^Version:.*/Version: ${VERSION}/" $SPEC_FILE

git-archive-all -C "$GIT_TOPDIR" --prefix="${PKG_NAME}-${VERSION}/" "$GIT_TOPDIR/rpm/${PKG_NAME}-${VERSION}.tar.gz"

go mod vendor
tar czf "vendor-${VERSION}.tar.gz" vendor/
mv "vendor-${VERSION}.tar.gz" "$GIT_TOPDIR/rpm/"
