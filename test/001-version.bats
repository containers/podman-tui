#!/usr/bin/env bats
#
# podman-tui version test
#

load helpers
load helpers_tui

@test "podman-tui version" {
    run_podman_tui version
    assert "$output" =~ "podman-tui v0.2.0-dev" "expected version"
}