#!/usr/bin/env bash
set -e

PKG_MANAGER=$(command -v dnf yum|head -n1)
${PKG_MANAGER} -y install \
    btrfs-progs-devel device-mapper-devel glib2-devel glibc-static \
    gpgme-devel libassuan-devel shadow-utils-subid-devel glibc-static gcc make golang \
    rpkg go-rpm-macros
