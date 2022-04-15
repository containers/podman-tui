## podman-tui

[![PkgGoDev](https://pkg.go.dev/badge/github.com/containers/podman-tui)](https://pkg.go.dev/github.com/containers/podman-tui)
[![Go Report](https://goreportcard.com/badge/github.com/containers/podman-tui)](https://goreportcard.com/report/github.com/containers/podman-tui)

podman-tui is a Terminal User Interface to interact with the podman v4.  
[podman bindings](https://github.com/containers/podman/tree/main/pkg/bindings) has been used to communicate with podman environment through rest api API (unix socket).

![Screenshot](./docs/podman-tui.gif)

---

## Installation

Building from source or installing packaged versions are detailed in [install guide](install.md).

---

## PreRun Checks

* podman-tui use podman unix socket for query therefore `podman.socket` service needs to be running.  
    The recommended way to start Podman system service in production mode is via systemd socket-activation:  

    ```shell
    $ systemctl --user start podman.socket
    ```

    See [start podman system service](https://podman.io/blogs/2020/08/10/podman-go-bindings.html) for more details.

* podman-tui uses 256 colors terminal mode. On `Nix system make sure TERM is set accordingly.

    ```shell
    $ export TERM=xterm-256color
    ```

---

## The Command Line
```shell
Usage:
  podman-tui [flags]
  podman-tui [command]

Available Commands:
  help        Help about any command
  version     Display podman-tui version and exit.


Flags:
  -d, --debug             Run application in debug mode
  -h, --help              help for podman-tui
  -l, --log-file string   Application runtime log file (default "podman-tui.log")

```

---

## Key Binding

podman-tui uses following keyboard keys for different actions:

| Action                           | Key        |
| -------------------------------- | ---------- |
| Display command menu             | m          |
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

---

## Available commands on different views

Check [podman-tui docs](./docs/README.md) for list of available commands on different pages (pods, containers, images, ...)


