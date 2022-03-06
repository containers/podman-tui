# Contributing To Podman TUI
We'd love your contribtion on the project!

## Topics
* [Developer Certificate of Origin (DCO)](#developer_certificate_of_origin)
* [Contributing to Podman TUI](#contributing-to-podman-tui)
* [Missing Features](#missing-features)

## Developer Certificate of Origin

Podman-tui enforces the Developer Certificate of Origin (DCO) on Pull Requests (PRs). This means that all commit messages must contain a signature line to indicate that the developer accepts the DCO.

Here is the full [text of the DCO][0], reformatted for readability:

    By making a contribution to this project, I certify that:

      (a) The contribution was created in whole or in part by me and I have the right to submit it under the open source 
          license indicated in the file; or

      (b) The contribution is based upon previous work that, to the best of my knowledge, is covered under an
          appropriate open source license and I have the right under that license to submit that work with modifications, 
          whether created in whole or in part by me, under the same open source license (unless I am permitted to submit 
          under a different license), as indicated in the file; or

      (c) The contribution was provided directly to me by some other person who certified (a), (b) or (c) and I have not 
          modified it.

      (d) I understand and agree that this project and the contribution are public and that a record of the contribution 
          (including all personal information I submit with it, including my sign-off) is maintained indefinitely and may 
          be redistributed consistent with this project or the open source license(s) involved.


Contributors indicate that they adhere to these requirements by adding
a `Signed-off-by` line to their commit messages.  For example:

    This is my commit message

    Signed-off-by: Random J Developer <random@developer.example.org>

The name and email address in this line must match those of the
committing author's GitHub account.

## Contributing To Podman TUI
### Fork and clone Podman TUI

First you need to fork and clone podman-tui project on Github.

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
  ```

#### Debian

  ```shell
  $ sudo apt-get -y install libgpgme-dev libbtrfs-dev libdevmapper-dev libassuan-dev pkg-config
  ```

### Validate, build binary and test your changes

```shell
$ make validate
$ make binary
$ ./bin/podman-tui
```

## Missing Features
```
* podman run
* podman system reset
* podman network connect
* podman network disconnect
* podman network reload
* podman container attach/detach
* remote connection
* cover more podman container create options 
* ... 

```
[0]: https://developercertificate.org/
