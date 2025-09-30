## Installation Guide

- [**Building From Source**](#building-from-source)
- [**Installing on Linux**](#installing-on-linux)
  - [**Alpine Linux**](#alpine-linux)
  - [**AlmaLinux, Rocky Linux**](#almalinux-rocky-linux)
  - [**Arch Linux (AUR)**](#arch-linux-aur)
  - [**CentOS Stream**](#centos-stream)
  - [**Fedora**](#fedora)
  - [**Gentoo**](#gentoo)
  - [**RHEL**](#rhel)
- [**Installing on Mac**](#installing-on-mac)
- [**Installing on Windows**](#installing-on-windows)
- [**Container Image**](#container-image)
- [**Configuration Files**](#configurations-files)

## Building From Source

podman-tui is using go version >= 1.20.

```shell
$ git clone <repository>
$ make binary # Linux
$ make binary-win # Windows
$ make binary-darwin # MacOS
```

## Installing on Linux

### Alpine Linux

```shell
$ sudo apk add podman-tui
```

### AlmaLinux, Rocky Linux

Enable [EPEL repository](https://docs.fedoraproject.org/en-US/epel/) and then run:

```shell
$ sudo dnf -y install podman-tui
```

### Arch Linux (AUR)

```shell
$ yay -S podman-tui
```

### CentOS Stream

Enable [EPEL repository](https://docs.fedoraproject.org/en-US/epel/) and then run:

```shell
$ sudo dnf -y install podman-tui
```

### Fedora

```shell
$ sudo dnf -y install podman-tui
```

### Gentoo

```shell
$ sudo emerge app-containers/podman-tui
```

### RHEL

Enable [EPEL repository](https://docs.fedoraproject.org/en-US/epel/) and then run:

```shell
$ sudo dnf -y install podman-tui
```

## Installing on Mac

podman-tui can be obtained through Homebrew package manager.

```shell
$ brew install podman-tui
```

# Installing on Windows

podman-tui can be obtained through WinGet Windows package manager.

`NOTE:` podman-tui is only supported on Windows PowerShell.

```shell
winget install Containers.PodmanTUI
```

## Container image

### Pull from quay.io

```shell
$ podman run -it --name podman-tui-app \
  -e CONTAINER_PASSPHRASE="<ssh key passphrase>" \
  -v <ssh_keys_dir>:/ssh_keys/:Z \
  --net=host \
  quay.io/navidys/podman-tui:latest # latest release, use develop tag to pull the upstream build
```

### Build image

podman-tui is using go version >= 1.17.

```shell
$ git clone <repository>
$ make binary
$ podman build -t podman-tui -f Containerfile
$ podman run -it --name podman-tui-app \
  -e CONTAINER_PASSPHRASE="<ssh key passphrase>" \
  -v <ssh_keys_dir>:/ssh_keys/:Z \
  --net=host \
  podman-tui
```


## Configuration Files

### podman-tui.json

~/.config/podman-tui/podman-tui.json

podman-tui.json is the configuration file which specifies local and remotes podman systems connections details.

```shell
{
  "connections": {
    "f42node01": {
      "uri": "ssh://navid@f42node01:22/run/user/1000/podman/podman.sock",
      "identity": "/home/navid/.ssh/id_ed25519"
    },
    "fc42node02": {
      "uri": "ssh://navid@f42node02:22/run/user/1000/podman/podman.sock",
      "identity": "/home/navid/.ssh/id_ed25519"
    },
    "localhost": {
      "uri": "unix://run/user/1000/podman/podman.sock",
      "default": true
    }
  }
}
```
