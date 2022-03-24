#!/usr/bin/env bats
#
# podman-tui system view functionality tests
#

load helpers
load helpers_tui

@test "system add" {
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
    sleep 1
    
    run_helper tail -2 $PODMAN_TUI_CONFIG_FILE
    assert "$output" =~ "[services.${TEST_SYSTEM_CONN_NAME}]" "expected [services.${TEST_SYSTEM_CONN_NAME}] in ${PODMAN_TUI_CONFIG_FILE}"
    assert "$output" =~ "uri = \"unix://run/podman/podman.sock\"" "expected ${TEST_SYSTEM_CONN_URI} in ${PODMAN_TUI_CONFIG_FILE}"
}

@test "system set default" {
    # switch to system view
    # select localhost_test connection name
    # select "set default" command
    podman_tui_set_view "system"
    podman_tui_select_item 1
    podman_tui_select_system_cmd "default"
    sleep 1
    
    run_helper tail -3 $PODMAN_TUI_CONFIG_FILE
    assert "$output" =~ "[services.${TEST_SYSTEM_CONN_NAME}]" "expected [services.${TEST_SYSTEM_CONN_NAME}] in ${PODMAN_TUI_CONFIG_FILE}"
    assert "$output" =~ "uri = \"unix://run/podman/podman.sock\"" "expected ${TEST_SYSTEM_CONN_URI} in ${PODMAN_TUI_CONFIG_FILE}"
    assert "$output" =~ "default = true" "expected 'default = true' in ${PODMAN_TUI_CONFIG_FILE}"
}

@test "system remove" {
    # switch to system view
    # select localhost_test connection name
    # select "remove connection" command
    # confirm connection removal
    podman_tui_set_view "system"
    podman_tui_select_item 1
    podman_tui_select_system_cmd "remove"
    podman_tui_send_inputs "Enter" 
    sleep 1
    
    run_helper tail -3 $PODMAN_TUI_CONFIG_FILE
    assert "$output" !~ "services.${TEST_SYSTEM_CONN_NAME}" "expected [services.${TEST_SYSTEM_CONN_NAME}] not in ${PODMAN_TUI_CONFIG_FILE}"
}

@test "system disconnect" {
    # switch to system view
    # select "disconnect" command
    podman_tui_set_view "system"
    podman_tui_select_system_cmd "disconnect"
    sleep 2

    run_helper tmux capture-pane -pS 0 -E 0
    assert "$output" =~ "DISCONNECTED" "expected DISCONNECTED connection status"

    run_helper tmux capture-pane -pS 7 -E 7
    assert "$output" !~ "connected" "expected empty connection status"
}

@test "system connect" {
    # switch to system view
    # select "disconnect" command
    podman_tui_set_view "system"
    podman_tui_select_system_cmd "disconnect"
    sleep 2
    run_helper tmux capture-pane -pS 0 -E 0
    assert "$output" =~ "DISCONNECTED" "expected DISCONNECTED connection status"

    # select "connect" command
    podman_tui_select_system_cmd "connect"
    sleep 3
    run_helper tmux capture-pane -pS 0 -E 0
    assert "$output" =~ "STATUS_OK" "expected STATUS_OK connection status"

    run_helper tmux capture-pane -pS 7 -E 7
    assert "$output" =~ "connected" "expected connected connection status"
}
