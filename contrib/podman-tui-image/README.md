# podman-tui-image

## Overview

This directory contains the Containerfile to create the podman-tui-image container
image.

The container image is built using the latest Fedora and then podman-tui
is built (upstream) and installed into it.

## Usage

```
# Build the podman-tui image
podman build -t podman-tui -f Containerfile

# Run the image and attach using the host's network
# Bind mount ssh keys volume for SSH connection to remote podman nodes
# Set SSH identity passphrase if required
podman run -it --name podman-tui-app -e CONTAINER_PASSPHRASE="<ssh key passphrase>" -v <ssh_keys_dir>:/ssh_keys/:Z --net=host podman-tui podman-tui
```
