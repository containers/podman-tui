#!/usr/bin/env bash
#
# Source file: https://raw.githubusercontent.com/containers/netavark/main/.packit.sh
#
# Packit will use the latest tag as version for spec file in case of rpkg.
# Using this script to update the spec file with the correct version.
#

set -eo pipefail

PKG_NAME="podman-tui"

# Get Version from HEAD
HEAD_VERSION=$(egrep "[[:space:]]+appVersion\s*=" cmd/version.go  | cut -d\" -f2 | sed -e 's/-dev//')

# Generate source tarball from HEAD
git archive --prefix=${PKG_NAME}-${HEAD_VERSION}/ -o ${PKG_NAME}-${HEAD_VERSION}.tar.gz HEAD

# RPM Spec modifications

# Update Version in spec with Version
sed -i "s/^Version:.*/Version: ${HEAD_VERSION}/" ${PKG_NAME}.spec

# Update Release in spec with Packit's release envvar
sed -i "s/^Release:.*/Release: $PACKIT_RPMSPEC_RELEASE%{?dist}/" ${PKG_NAME}.spec

# Update Source tarball name in spec
sed -i "s/^Source0:.*.tar.gz/Source0: %{name}-${HEAD_VERSION}.tar.gz/" ${PKG_NAME}.spec

# Update setup macro to use the correct build dir
sed -i "s/^%setup.*/%setup -T -b 0 -q -n %{name}-${HEAD_VERSION}/" ${PKG_NAME}.spec
