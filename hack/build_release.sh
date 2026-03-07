#!/bin/bash
SRC_DIR=$(dirname `realpath $0`)
VERSION="${1#v}"
BUILD_DIR="${SRC_DIR}/../dist/"
CUR_DIR=$(pwd)

mkdir -p ${BUILD_DIR}/podman-tui-release-linux_amd64/podman-tui-${VERSION}/
mkdir -p ${BUILD_DIR}/podman-tui-release-darwin_amd64/podman-tui-${VERSION}/
mkdir -p ${BUILD_DIR}/podman-tui-release-windows_amd64/podman-tui-${VERSION}/

make all

cp -p ./bin/podman-tui ${BUILD_DIR}/podman-tui-release-linux_amd64/podman-tui-${VERSION}/
cp -p ./bin/darwin/podman-tui ${BUILD_DIR}/podman-tui-release-darwin_amd64/podman-tui-${VERSION}/
cp -p ./bin/windows/podman-tui.exe ${BUILD_DIR}/podman-tui-release-windows_amd64/podman-tui-${VERSION}/

echo "############ linux amd64 ############"
file ${BUILD_DIR}/podman-tui-release-linux_amd64/podman-tui-${VERSION}/podman-tui
${BUILD_DIR}/podman-tui-release-linux_amd64/podman-tui-${VERSION}/podman-tui version

echo "############ darwin amd64 ############"
file ${BUILD_DIR}/podman-tui-release-darwin_amd64/podman-tui-${VERSION}/podman-tui

echo "############ windows amd64 ############"
file ${BUILD_DIR}/podman-tui-release-windows_amd64/podman-tui-${VERSION}/podman-tui.exe


cd ${BUILD_DIR}
zip -r podman-tui-release-linux_amd64.zip podman-tui-release-linux_amd64
zip -r podman-tui-release-darwin_amd64.zip podman-tui-release-darwin_amd64
zip -r podman-tui-release-windows_amd64.zip podman-tui-release-windows_amd64
sha256sum *.zip > sha256sum
cd ${CUR_DIR}