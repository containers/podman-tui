## podman-tui

![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/containers/podman-tui)
![Go](https://github.com/containers/podman-tui/workflows/Go/badge.svg)
[![Go Report](https://goreportcard.com/badge/github.com/containers/podman-tui)](https://goreportcard.com/report/github.com/containers/podman-tui)
[![codecov](https://codecov.io/gh/navidys/podman-tui/branch/main/graph/badge.svg)](https://codecov.io/gh/navidys/podman-tui)

Terminal user interface for Podman environment.

![Screenshot](./podman-tui.gif)

- [**Overview**](#overview)
- [**Compatibility Matrix**](#compatibility-matrix)
- [**Installation**](#installation)
- [**PreRun Checks**](#prerun-checks)
- [**Key Bindings**](#key-bindings)
- [**Code of Conduct**](#code-of-conduct)
- [**License**](#license)

## Overview

podman-tui is a terminal user interface for podman environment.
It is using [podman go bindings](https://github.com/containers/podman/tree/main/pkg/bindings) to communicate with local or remote podman machine (through SSH).

## Compatibility matrix

| Terminal User Interface  | Podman |
| ----------- | ------ |
| release-1.x | v5.x.y |
| release-0.x | v4.x.y |

## Installation

Building from source (Linux, Windows and MacOS) or installing packaged versions are detailed in [install guide](install.md).

## PreRun Checks

* `podman.socket` service needs to be running on podman machine.
    The recommended way to start Podman system service in production mode is via systemd socket-activation:

    ```shell
    $ systemctl --user start podman.socket
    ```

    For non-systemd distributions, the socket can be created using the below command.
    Add this command to your distro or window manager’s autostart to start the service automatically.

    ```shell
    $ podman system service --time=0
    ```

    Note: This command creates the socket only for the user who runs it.

    See [start podman system service](https://podman.io/blogs/2020/08/10/podman-go-bindings.html) for more details.

* If the SSH key has a passphrase, then you need to set and export `CONTAINER_PASSPHRASE=<password>` variable.
    ```shell
    $ export CONTAINER_PASSPHRASE=keypass
    ```

* podman-tui uses 256 colors terminal mode on `Nix` system and standard ANSI colors on `Windows` system (use "command prompt" application).


## Key Bindings

podman-tui uses following keyboard keys for different actions:

| Action                           | Key        |
| -------------------------------- | ---------- |
| Display command menu             | m          |
| Display sort menu                | s          |
| Switch to next screen            | l          |
| Switch to previous screen        | h          |
| Move up                          | k          |
| Move down                        | j          |
| Exit application                 | Ctrl+c     |
| Close the active dialog          | Esc        |
| Switch between interface widgets | Tab        |
| Delete selected item             | Delete     |
| Move up/down                     | Up/Down    |
| Previous/Next screen             | Left/Right |
| Scroll Up                        | Page Up    |
| Scroll Down                      | Page Down  |
| Display help screen              | F1         |
| Display system screen            | F2         |
| Display pods screen              | F3         |
| Display containers screen        | F4         |
| Display volumes screen           | F5         |
| Display images screen            | F6         |
| Display networks screen          | F7         |
| Display secrets screen           | F8         |

## Code of Conduct

This project is using the [Containers Community Code of Conduct](https://github.com/containers/common/blob/main/CODE-OF-CONDUCT.md)

## License

Licensed under the [Apache 2.0](LICENSE) license.
