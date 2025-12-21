#!/usr/bin/env bats
#
# podman-tui networks view functionality tests
#

load helpers
load helpers_tui

@test "network connect" {
    check_skip "network_connect"

    podman network rm $TEST_NETWORK_CONNECT || echo done
    podman container rm -f $TEST_CONTAINER_NAME || echo done
    podman container create --name $TEST_CONTAINER_NAME docker.io/library/busybox || echo done
    podman network create $TEST_NETWORK_CONNECT || echo done
    # switch to networks view
    # select connect command from network commands dialog
    # select container
    # select connect button

    podman_tui_set_view "networks"
    podman_tui_select_item 1
    podman_tui_select_network_cmd "connect"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab"
    podman_tui_send_inputs $TEST_NETWORK_CONNECT_ALIAS
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Tab" "Enter"

    sleep $TEST_TIMEOUT_LOW

    run_helper podman container inspect $TEST_CONTAINER_NAME  --format "\"{{ .NetworkSettings.Networks.$TEST_NETWORK_CONNECT }}\""
    assert "$output" =~ "$TEST_NETWORK_CONNECT_ALIAS" "expected $TEST_NETWORK_CONNECT_ALIAS to be in the list of aliases"

}

@test "network disconnect" {
    check_skip "network_disconnect"

    # switch to networks view
    # select disconnect command from network commands dialog
    # select container
    # select disconnect button

    podman_tui_set_view "networks"
    podman_tui_select_item 1
    podman_tui_select_network_cmd "disconnect"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab" "Tab" "Enter"

    run_helper podman container inspect $TEST_CONTAINER_NAME  --format "{{ .NetworkSettings.Networks.$TEST_NETWORK_CONNECT }}"
    assert "$output" == "<no value>" "expected $TEST_NETWORK_CONNECT_ALIAS to be removed from container"

}

@test "network create" {
    check_skip "network_create"

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
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Tab" "Enter"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab" "Enter"
    sleep $TEST_TIMEOUT_LOW
    run_helper podman network ls --format "{{ .Name }}" --filter "name=${TEST_NETWORK_NAME}$"
    assert "$output" == "$TEST_NETWORK_NAME" "expected $TEST_NETWORK_NAME to be in the list"
}

@test "network inspect" {
    check_skip "network_inspect"

    # switch to networks view
    # select test network from list
    # select inspect command from network commands dialog
    # close network inspect result message dialog
    podman_tui_set_view "networks"
    podman_tui_select_item 1
    podman_tui_select_network_cmd "inspect"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper sed -n '/  "labels": {/, /  }/p' $PODMAN_TUI_LOG
    assert "$output" =~ "\"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\"" "expected \"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\" in network inspect"
}

@test "network remove" {
    check_skip "network_remove"

    # switch to networks view
    # select test network from list
    # select remove command from network commands dialog
    podman_tui_set_view "networks"
    podman_tui_select_item 1
    podman_tui_select_network_cmd "remove"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman network ls --format "{{ .Name }}" --filter "name=${TEST_NETWORK_NAME}$"
    assert "$output" == "" "expected $TEST_NETWORK_NAME removed"

}

@test "network prune" {
    check_skip "network_prune"

    run_helper podman network create $TEST_NETWORK_NAME

    # switch to networks view
    # select prune command from network commands dialog
    # confirm the operation on warnings dialog
    podman_tui_set_view "networks"
    podman_tui_select_network_cmd "prune"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman network ls --format "{{ .Name }}" --filter "name=${TEST_NETWORK_NAME}$"
    assert "$output" == "" "expected at least $TEST_NETWORK_NAME network removal"
}
