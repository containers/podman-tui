#!/bin/bash

# Description:
# This script creates an archive with vendored dependencies from a Go SPEC file.

# License:
# MIT License
#
# Copyright (c) 2023 Robert-André Mauchin <zebob.m@gmail.com>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# Check if the RPM SPEC file is given as an argument
if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <path_to_rpm_spec_file>"
    exit 1
fi

RPM_SPEC_FILE=$1

# Extract the directory from the RPM SPEC file path
SPEC_DIR=$(dirname $(realpath "$RPM_SPEC_FILE"))

# Extract the URL, commit, tag, and version from the RPM SPEC file
FORGEURL=$(awk '/^%global forgeurl/ {print $NF}' "$RPM_SPEC_FILE")
GOIPATH=$(awk '/^%global goipath/ {print $NF}' "$RPM_SPEC_FILE")
COMMIT=$(awk '/^%global commit/ {print $NF}' "$RPM_SPEC_FILE")
TAG=$(awk '/^%global tag/ {print $NF}' "$RPM_SPEC_FILE")
VERSION=$(awk '/^Version:/ {print $NF}' "$RPM_SPEC_FILE")

# Decide which URL to use
if [[ -n "$FORGEURL" ]]; then
    REPO_URL="$FORGEURL"
elif [[ -n "$GOIPATH" ]]; then
    REPO_URL="https://$GOIPATH"
else
    echo "No repository URL found in the RPM SPEC file."
    exit 2
fi

# Create a temporary directory and clone the repository
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT
git clone "$REPO_URL" "$TMP_DIR"
if [[ $? -ne 0 ]]; then
    echo "Failed to clone repository."
    exit 3
fi

# Change to the directory
pushd "$TMP_DIR" > /dev/null

# Checkout based on priority: commit > tag > Version
CHECKOUT_SUCCESS=0
if [[ -n "$COMMIT" ]]; then
    CHECKOUT_IDENTIFIER="$COMMIT"
    git checkout "$CHECKOUT_IDENTIFIER" && CHECKOUT_SUCCESS=1
elif [[ -n "$TAG" ]]; then
    CHECKOUT_IDENTIFIER="$TAG"
    git checkout "$CHECKOUT_IDENTIFIER" && CHECKOUT_SUCCESS=1
elif [[ -n "$VERSION" ]]; then
    CHECKOUT_IDENTIFIER="$VERSION"
    git checkout "$VERSION" || git checkout "v$VERSION"
    if [ $? -eq 0 ]; then
        CHECKOUT_SUCCESS=1
    fi
else
    echo "No commit, tag, or version found in the RPM SPEC file."
    exit 4
fi

if [ $CHECKOUT_SUCCESS -eq 0 ]; then
    echo "Failed to checkout using commit, tag, or version."
    exit 5
fi

# Run go mod vendor
go mod vendor
if [[ $? -ne 0 ]]; then
    echo "Failed to run 'go mod vendor'."
    exit 6s
fi

# Create a tar.gz of the vendor directory
tar czf "vendor-$CHECKOUT_IDENTIFIER.tar.gz" vendor/
if [[ $? -ne 0 ]]; then
    echo "Failed to create tar.gz of the vendor directory."
    exit 7
fi

# Move the tar.gz to the SPEC directory
mv "vendor-$CHECKOUT_IDENTIFIER.tar.gz" "$SPEC_DIR/"

# Go back to the original directory
popd > /dev/null

# Clean up
rm -rf "$TMP_DIR"

echo "Created vendor-$CHECKOUT_IDENTIFIER.tar.gz successfully."
