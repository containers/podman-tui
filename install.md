## Installation Guide

- [**Building From Source**](#building-from-source)
- [**Installing Packaged Versions**](#installing-packaged-versions)
  - [**Arch Linux (AUR)**](#arch-linux-aur)
  - [**Fedora**](#fedora)
  - [**Gentoo**](#gentoo)
- [**Container Image**](#container-image)
- [**Configuration Files**](#configurations-files)

## Building From Source

podman-tui is using go version >= 1.17.

```shell
$ git clone <repository>
$ make binary # Linux
$ make binary-win # Windows
$ make binary-darwin # MacOS
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

### Gentoo

```
$ sudo emerge app-containers/podman-tui
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

### podman-tui.conf

~/.config/podman-tui/podman-tui.conf

podman-tui.conf is the configuration file which specifies local and remotes podman systems connections details.

```shell
[services]

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
