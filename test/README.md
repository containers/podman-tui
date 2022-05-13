# Podman-tui functionality tests with bats

## Running tests

To run the tests locally in your sandbox, you can use one of these methods:

```shell
$ make binary                    # build podman-tui binary
$ sudo make test                 # run all the tests
$ sudo bats test/001-image.bats  # runs just the specified test
```

## Requirements
- switch user to `root` account before running the tests
- use provided Vagrant file to create clean VM sandbox (optional).
- if you are not in root directory of the project, be sure `PODMAN_TUI` variable is set
- access to repository to pull busybox and httpd image
- tmux
