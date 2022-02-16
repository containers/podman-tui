# Podman-tui functionality tests with bats


## Running tests

To run the tests locally in your sandbox, you can use one of these methods:

* bats ./test/001-image.bats  # runs just the specified test
* bats ./test/                # runs all

## Requirements
- busybox image shall not exist on the sandbox (it will be pulled automatically during the test)
- tmux