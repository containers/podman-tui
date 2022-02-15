#!/usr/bin/env bats
#
# podman-tui networks view functionality tests
#

load helpers
load helpers_tui

@test "network create" {
    podman network rm $TEST_NETWORK_NAME || echo done
    
    # switch to networks view
    # select create command from network commands dialog
    # fillout create dialog fields and press enter
    # close network create result message dialog
    podman_tui_set_view "networks"
    podman_tui_select_network_cmd "create"
    podman_tui_send_inputs "$TEST_NETWORK_NAME"
    podman_tui_send_inputs "Tab"
    podman_tui_send_inputs "$TEST_LABEL"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Enter"
    sleep 2
    podman_tui_send_inputs "Tab" "Enter"

    run_podman network ls --format "'{{ .Name }}'" --filter "name=$TEST_NETWORK_NAME"
    assert "$output" =~ "$TEST_NETWORK_NAME" "expected $TEST_NETWORK_NAME to be in the list"
}

@test "network inspect" {
    net_index=$(podman network ls -q | nl -v 0 | grep "$TEST_NETWORK_NAME" | awk '{print $1}')
    
    # switch to networks view
    # select test network from list
    # select inspect command from network commands dialog
    # close network inspect result message dialog
    podman_tui_set_view "networks"
    podman_tui_select_item $net_index
    podman_tui_select_network_cmd "inspect"
    sleep 2
    podman_tui_send_inputs "Enter"

    run_helper sed -n '/  "podman_labels": {/, /  }/p' $PODMAN_TUI_LOG
    assert "$output" "=~" "\"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\"" "expected \"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\" in network inspect"
}

@test "network remove" {
    net_index=$(podman network ls -q | nl -v 0 | grep "$TEST_NETWORK_NAME" | awk '{print $1}')

    # switch to networks view
    # select test network from list
    # select remove command from network commands dialog
    podman_tui_set_view "networks"
    podman_tui_select_item $net_index
    podman_tui_select_network_cmd "remove"
    podman_tui_send_inputs "Enter"
    sleep 2

    run_podman network ls --format "'{{ .Name }}'" --filter "name=$TEST_NETWORK_NAME"
    assert "$output" =~ "" "expected $TEST_NETWORK_NAME removed"

}

@test "network prune" {
    run_podman network create $TEST_NETWORK_NAME

    # switch to networks view
    # select prune command from network commands dialog
    # confirm the operation on warnings dialog
    podman_tui_set_view "networks"
    podman_tui_select_volume_cmd "prune"
    podman_tui_send_inputs "Enter"
    sleep 2

    run_podman network ls --format "'{{ .Name }}'" --filter "name=$TEST_NETWORK_NAME"
    assert "$output" "=~" "" "expected at least $TEST_NETWORK_NAME network removal"
}
