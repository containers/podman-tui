# Podman-tui functionality tests with bats

## Running tests

To run the tests locally in your sandbox, you can use one of these methods:

```shell
$ make binary                                # build podman-tui binary
$ sudo podman system reset                   # to reset all previous configuration
$ sudo systemctl restart podman.socket       # restart podman socket
$ sudo make test-functionality               # run all the tests
```

## Requirements
- use provided Vagrant file to create clean VM sandbox (optional).
- if you are not in root directory of the project, be sure `PODMAN_TUI` variable is set
- access to repository to pull busybox and httpd image
- tmux
