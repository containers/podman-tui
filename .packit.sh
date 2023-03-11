#!/usr/bin/env bash

# Used from https://raw.githubusercontent.com/containers/netavark/main/.packit.sh
# Packit's default fix-spec-file often doesn't fetch version string correctly.
# This script handles any custom processing of the dist-git spec file and gets used by the
# fix-spec-file action in .packit.yaml

set -eo pipefail

# Get Version from HEAD
HEAD_VERSION=$(cat cmd/version.go  | egrep "appVersion\s*=" | awk -F= '{print $2}' | sed "s/\"//g" | sed "s/ //g" | sed -e 's/-/~/')

# Generate source tarball from HEAD
git archive --prefix=podman-tui-$HEAD_VERSION/ -o podman-tui-$HEAD_VERSION.tar.gz HEAD

# RPM Spec modifications

# Update Version in spec with Version
sed -i "s/^Version:.*/Version: $HEAD_VERSION/" podman-tui.spec

# Update Release in spec with Packit's release envvar
sed -i "s/^Release:.*/Release: $PACKIT_RPMSPEC_RELEASE%{?dist}/" podman-tui.spec

# Update Source tarball name in spec
sed -i "s/^Source:.*.tar.gz/Source: %{name}-$HEAD_VERSION.tar.gz/" podman-tui.spec

# Update setup macro to use the correct build dir
sed -i "s/^%setup.*/%autosetup -Sgit -n %{name}-$HEAD_VERSION/" podman-tui.spec
