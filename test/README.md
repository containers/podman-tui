# Podman-tui functionality tests with bats

## Running tests

To run the tests locally in your sandbox, you can use one of these methods:

```shell
$ sudo /bin/su -
$ make binary          # build podman-tui binary
$ make test            # run all the tests
```

```shell
$ sudo /bin/su -
$ make binary          # build podman-tui binary
$ export PODMAN_TUI=$(realpath ./bin/podman-tui)
$ cd test/
$ bats 001-image.bats  # runs just the specified test
$ bats *               # run all
```

## Requirements
- switch user to `root` account before running the tests
- if you are not in root directory of the project, be sure `PODMAN_TUI` variable is set
- access to repository to pull busybox and httpd image 
- tmux