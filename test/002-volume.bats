#!/usr/bin/env bats
#
# podman-tui volumes view functionality tests
#

load helpers
load helpers_tui

@test "volume create" {
    check_skip "volume_create"

    podman volume rm $TEST_VOLUME_NAME || echo done

    # switch to volumes view
    # select create command from volume commands dialog
    # fillout create dialog fields and press enter
    # close volume create result message dialog
    podman_tui_set_view "volumes"
    podman_tui_select_volume_cmd "create"
    podman_tui_send_inputs "$TEST_VOLUME_NAME" "Tab" "$TEST_LABEL" "Tab" "Tab" "Tab" "Tab" "Enter"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab" "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman volume ls --format "{{ .Name }}" --filter "name=${TEST_VOLUME_NAME}"
    assert "$output" == "$TEST_VOLUME_NAME" "expected $TEST_VOLUME_NAME to be in the list"
}

@test "volume inspect" {
    check_skip "volume_inspect"

    # switch to volumes view
    # select test volume from list
    # select inspect command from volume commands dialog
    # close volume inspect result message dialog
    podman_tui_set_view "volumes"
    podman_tui_select_item 0
    podman_tui_select_volume_cmd "inspect"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper sed -n '/  "Labels": {/, /  },/p' ${PODMAN_TUI_LOG}
    assert "$output" =~ "\"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\"" "expected \"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\" in volume inspect"
}

@test "volume remove" {
    check_skip "volume_remove"

    # switch to volumes view
    # select test volume from list
    # select remove command from volume commands dialog
    podman_tui_set_view "volumes"
    podman_tui_select_item 0
    podman_tui_select_volume_cmd "remove"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman volume ls --format "{{ .Name }}" --filter "name=${TEST_VOLUME_NAME}"
    assert "$output" == "" "expected $TEST_VOLUME_NAME removed"

}

@test "volume prune" {
    check_skip "volume_prune"

    run_helper podman volume create $TEST_VOLUME_NAME

    # switch to volumes view
    # select prune volume from volume commands dialog
    # confirm the operation on warnings dialog
    podman_tui_set_view "volumes"
    podman_tui_select_volume_cmd "prune"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman volume ls --format "{{ .Name }}" --filter "name=${TEST_NETWORK_NAME}"
    assert "$output" =~ "" "expected at least $TEST_VOLUME_NAME image removal"
}
