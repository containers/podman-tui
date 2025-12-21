#!/usr/bin/env bats
#
# podman-tui system view functionality tests
#

load helpers
load helpers_tui

@test "system add" {
    check_skip "system_add"

    # switch to system view
    # select add connection
    # fillout name field
    # fillout URI field
    # go to Add button and press Enter
    podman_tui_set_view "system"
    podman_tui_select_system_cmd "add"
    podman_tui_send_inputs $TEST_SYSTEM_CONN_NAME
    podman_tui_send_inputs "Tab"
    podman_tui_send_inputs $TEST_SYSTEM_CONN_URI
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper jq ".connections.${TEST_SYSTEM_CONN_NAME}.uri" ${PODMAN_TUI_CONFIG_FILE}
    assert "$output" == "\"$TEST_SYSTEM_CONN_URI\"" "expected ${TEST_SYSTEM_CONN_URI} in ${PODMAN_TUI_CONFIG_FILE}"
}

@test "system set default" {
    check_skip "system_default"

    # switch to system view
    # select localhost_test connection name
    # select "set default" command
    podman_tui_set_view "system"
    podman_tui_select_item 1
    podman_tui_select_system_cmd "default"
    sleep $TEST_TIMEOUT_LOW

    run_helper jq ".connections.${TEST_SYSTEN_CONN_LOCAL}.uri" ${PODMAN_TUI_CONFIG_FILE}
    assert "$output" == "\"$TEST_SYSTEM_CONN_URI\"" "expected ${TEST_SYSTEM_CONN_URI} in ${PODMAN_TUI_CONFIG_FILE}"

    run_helper jq ".connections.${TEST_SYSTEN_CONN_LOCAL}.default" ${PODMAN_TUI_CONFIG_FILE}
    assert "$output" == "true" "expected 'default = true' in ${PODMAN_TUI_CONFIG_FILE}"
}

@test "system remove" {
    check_skip "system_remove"

    # switch to system view
    # select localhost_test connection name
    # select "remove connection" command
    # confirm connection removal
    podman_tui_set_view "system"
    podman_tui_select_item 1
    podman_tui_select_system_cmd "remove"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper jq ".connections.${TEST_SYSTEN_CONN_LOCAL}" ${PODMAN_TUI_CONFIG_FILE}
    assert "$output" == "null" "expected ${TEST_SYSTEN_CONN_LOCAL} connection to be removed from in ${PODMAN_TUI_CONFIG_FILE}"
}

@test "system disconnect" {
    check_skip "system_disconnect"

    # switch to system view
    # select "disconnect" command
    podman_tui_set_view "system"
    podman_tui_select_system_cmd "disconnect"
    sleep $TEST_TIMEOUT_LOW

    run_helper tmux capture-pane -pS 0 -E 0
    assert "$output" =~ "DISCONNECTED" "expected DISCONNECTED connection status"

    run_helper tmux capture-pane -pS 7 -E 7
    assert "$output" !~ "connected" "expected empty connection status"
}

@test "system connect" {
    check_skip "system_connect"

    # switch to system view
    # select "disconnect" command
    podman_tui_set_view "system"
    podman_tui_select_system_cmd "disconnect"
    sleep $TEST_TIMEOUT_LOW
    run_helper tmux capture-pane -pS 0 -E 0
    assert "$output" =~ "DISCONNECTED" "expected DISCONNECTED connection status"

    # select "connect" command
    podman_tui_select_system_cmd "connect"
    sleep $TEST_TIMEOUT_LOW
    run_helper tmux capture-pane -pS 0 -E 0
    assert "$output" =~ "STATUS_OK" "expected STATUS_OK connection status"

    run_helper tmux capture-pane -pS 7 -E 7
    assert "$output" =~ "connected" "expected connected connection status"
}
