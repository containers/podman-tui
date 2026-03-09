#!/usr/bin/env bash
PKG_NAME="podman-tui"
GIT_TOPDIR=$(git rev-parse --show-toplevel)
GIT_LAST_COMMIT=$(git log -1 --format=%H)
SPEC_FILE=${GIT_TOPDIR}/rpm/${PKG_NAME}.spec

VERSION=$(grep -E "[[:space:]]+appVersion\s*=" ${GIT_TOPDIR}/cmd/version.go  | cut -d\" -f2 | sed 's/-//')
RELEASE_VERSION=1

echo "$VERSION" | grep 'dev' 2> /dev/null
if [ $? -eq 0 ] ; then
    RELEASE_VERSION=0
fi

chmod +x ${GIT_TOPDIR}/rpm/bundle_go_deps_for_rpm.sh

sed -i "s/^Version:.*/Version: ${VERSION}/" $SPEC_FILE

git-archive-all -C "$GIT_TOPDIR" --prefix="${PKG_NAME}-${VERSION}/" "$GIT_TOPDIR/rpm/${PKG_NAME}-${VERSION}.tar.gz"

if [ $RELEASE_VERSION -eq 0 ] ; then
    go mod vendor
    tar czf "vendor-${VERSION}.tar.gz" vendor/
    mv "vendor-${VERSION}.tar.gz" "$GIT_TOPDIR/rpm/"
else
    ${GIT_TOPDIR}/rpm/bundle_go_deps_for_rpm.sh ${SPEC_FILE}
fi
