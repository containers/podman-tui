# Contributing To Podman TUI
We'd love your contribtion on the project!

## Topics
* [Contributing to Podman TUI](#contributing-to-podman-tui)
* [Missing Features](#missing-features)

## Contributing To Podman TUI

### Fork and clone Podman TUI

First you need to frok and clone podman-tui project on Github.

Be sure to have [defined your `$GOPATH` environment variable](https://github.com/golang/go/wiki/GOPATH).

Create a path that corresponds to the go import paths of Podman-tui: `mkdir -p $GOPATH/src/github.com/containers`.

Then clone your fork locally:

    ```shell
    $ git clone git@github.com:<you>/podman-tui $GOPATH/src/github.com/containers/podman-tui
    $ cd $GOPATH/src/github.com/containers/podman-tui
    ```

### Deal with make

Podman TUI use a Makefile to realize common action like building etc...

You can list available actions by using:

```shell
$ make help
Usage: make <target>
...output...
```

### Prerequisite before build
You need install some dependencies before building a binary.

#### Fedora

  ```shell
  $ sudo dnf install -y btrfs-progs-devel device-mapper-devel gpgme-devel libassuan-devel
  $ export PKG_CONFIG_PATH="/usr/lib/pkgconfig"
  ```

### Building binaries and test your changes

To test you changes do `make binary` to generate your binary.

Your binary is created in ise the bin/ directory and you can test your changes:

```shell
$ bin/podman-tui
```

## Missing Features
```
* podman exec
* podman system reset
* podman network connect
* podman network disconnect
* podman newtork reload
* podman pod stats
* podman contaienrs stats
* remote connection
* cover more podman container create options 
* ... 

```
