## Installation Guide

- [**Building From Source**](#building-from-source)
- [**Installing Packaged Versions**](#installing-packaged-versions)
  - [**Arch Linux (AUR)**](#arch-linux-aur)
  - [**Fedora**](#fedora)
- [**Configuration Files**](#configurations-files)


## Building From Source

podman-tui is using go version >= 1.17. 
  1. Clone the repo
  2. Install [dependencies](./CONTRIBUTING.md#prerequisite-before-build)
  3. Build binaries (Linux and Windows)
     
     ```shell
     $ make all
     ```

## Installing Packaged Versions

### Arch Linux (AUR)

```shell
$ yay -S podman-tui
```

### Fedora

```
$ sudo dnf -y install podman-tui
```

## Configuration Files

### podman-tui.conf

~/.config/podman-tui/podman-tui.conf

podman-tui.conf is the configuration file which specifies local and remotes podman systems connections details.

```shell
services]

  [services.fc36node01]
    uri = "ssh://navid@fc36node01:22/run/user/1000/podman/podman.sock"
    identity = "/home/navid/.ssh/id_ed25519"
  [services.fc36node02]
    uri = "ssh://navid@fc36node02:22/run/user/1000/podman/podman.sock"
    identity = "/home/navid/.ssh/id_ed25519"
    default = true
  [services.localhost]
    uri = "unix://run/user/1000/podman/podman.sock"
```