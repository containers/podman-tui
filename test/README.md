# Podman-tui functionality tests with bats


## Running tests

To run the tests locally in your sandbox, you can use one of these methods:

```shell
$ make binary          # build podman-tui binary
$ export PODMAN_TUI=$(realpath ./bin/podman-tui)
$ cd ./tests/
$ bats 001-image.bats  # runs just the specified test
$ bats *               # run all
```

## Requirements
- if you are not on root directory of the project, be sure `PODMAN_TUI` variable is set
- access to repository to pull busybox and httpd image 
- tmux