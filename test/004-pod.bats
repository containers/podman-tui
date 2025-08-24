#!/usr/bin/env bats
#
# podman-tui pods view functionality tests
#

load helpers
load helpers_tui

@test "pod create (resource)" {
    podman pod rm -f $TEST_POD_NAME || echo done
    podman image pull pause:3.5 || echo done

    podman_tui_set_view "pods"
    podman_tui_select_pod_cmd "create"
    podman_tui_send_inputs $TEST_POD_NAME "Tab" "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Down" "Down" "Down" "Down" "Down" "Tab"
    podman_tui_send_inputs $TEST_POD_MEMORY "Tab" $TEST_POD_SWAP
    podman_tui_send_inputs "Tab" "Tab" $TEST_POD_CPU_SHARES
    podman_tui_send_inputs "Tab" "Tab" $TEST_POD_CPUSET_MEM
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab"
    sleep $TEST_TIMEOUT_LOW

    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --filter="name=${TEST_POD_NAME}$" --format "{{ .Status}}"
    assert $output =~ "Created" "expected $TEST_POD_NAME to be created"

    pod_memory=$(podman pod inspect $TEST_POD_NAME --format "{{ json .MemoryLimit }}")
    pod_memory_swap=$(podman pod inspect $TEST_POD_NAME --format "{{ json .MemorySwap }}")
    pod_cpu_shares=$(podman pod inspect $TEST_POD_NAME --format "{{ json .CPUShares }}")
    pod_cpu_set_mems=$(podman pod inspect $TEST_POD_NAME --format "{{ json .CPUSetMems }}")

    assert "$pod_memory" =~ "$TEST_POD_MEMORY" "expected pod memory to match: $TEST_POD_MEMORY"
    assert "$pod_memory_swap" =~ "$TEST_POD_SWAP" "expected pod memory swap to match: $TEST_POD_SWAP"
    assert "$pod_cpu_shares" =~ "$TEST_POD_CPU_SHARES" "expected pod cpu shares to match: $TEST_POD_CPU_SHARES"
    assert "$pod_cpu_set_mems" =~ "$TEST_POD_CPUSET_MEM" "expected pod cpu set mems to match: $TEST_POD_CPUSET_MEM"
}

@test "pod create (networking, security)" {
    podman pod rm -f $TEST_POD_NAME || echo done
    podman network rm $TEST_POD_NETWORK_NAME || echo done
    podman image pull pause:3.5 || echo done
    podman network create $TEST_POD_NETWORK_NAME || echo done

    # switch to pods view
    # select create command from pod commands dialog
    # fillout name field
    # fillout label field
    # switch to "networking" create view
    # go to networks dropdown widget and select network name from available networks
    # switch to "security option" create view
    # set label disable
    # set no new privileges
    # go to "Create" button and press Enter
    podman_tui_set_view "pods"
    podman_tui_select_pod_cmd "create"
    podman_tui_send_inputs $TEST_POD_NAME "Tab" "Tab" $TEST_LABEL
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Down" "Down" "Down" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Down"
    podman_tui_select_item 1
    podman_tui_send_inputs "Enter"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Down" "Tab"
    podman_tui_send_inputs "disable"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Tab" "Space"
    podman_tui_send_inputs "Tab" "Tab"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --filter="name=${TEST_POD_NAME}$" --format "{{ .Status}}"
    assert $output =~ "Created" "expected $TEST_POD_NAME to be created"

    security_opts=$(podman pod inspect $TEST_POD_NAME | sed -n '/security_opt/,/.*]/p')
    assert "$security_opts" =~ "no-new-privileges" "expected no-new-privileges in pod security options"
    assert "$security_opts" =~ "label=disable" "expected label=disable in pod security options"
}

@test "pod start" {
    # switch to pods view
    # select test pod from list
    # select start command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item 0
    podman_tui_select_pod_cmd "start"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --filter="name=${TEST_POD_NAME}$" --format "{{ .Status}}"
    assert $output =~ "Running" "expected $TEST_POD_NAME running"
}

@test "pod pause" {
    # switch to pods view
    # select test pod from list
    # select pause command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item 0
    podman_tui_select_pod_cmd "pause"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --filter="name=${TEST_POD_NAME}$" --format "{{ .Status}}"
    assert $output =~ "Paused" "expected $TEST_POD_NAME running"
}

@test "pod unpause" {
    # switch to pods view
    # select test pod from list
    # select unpause command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item 0
    podman_tui_select_pod_cmd "unpause"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --filter="name=${TEST_POD_NAME}$" --format "{{ .Status}}"
    assert $output =~ "Running" "expected $TEST_POD_NAME running"
}

@test "pod stop" {
    # switch to pods view
    # select test pod from list
    # select stop command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item 0
    podman_tui_select_pod_cmd "stop"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --filter="name=${TEST_POD_NAME}$" --format "{{ .Status}}"
    assert $output =~ "Exited" "expected $TEST_POD_NAME exited"
}

@test "pod restart" {
    # switch to pods view
    # select test pod from list
    # select restart command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item 0
    podman_tui_select_pod_cmd "restart"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --filter="name=${TEST_POD_NAME}$" --format "{{ .Status}}"
    assert $output =~ "Running" "expected $TEST_POD_NAME exited"
}

@test "pod kill" {
    # switch to pods view
    # select test pod from list
    # select kill command from pod commands dialog
    podman_tui_set_view "pods"
    podman_tui_select_item 0
    podman_tui_select_pod_cmd "kill"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --filter="name=${TEST_POD_NAME}$" --format "{{ .Status}}"
    assert $output =~ "Exited" "expected $TEST_POD_NAME exited"
}

@test "pod inspect" {
    # switch to pods view
    # select test pod from list
    # select inspect command from pod commands dialog
    # close pod inspect result message dialog
    podman_tui_set_view "pods"
    podman_tui_select_item 0
    podman_tui_select_pod_cmd "inspect"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Enter"

    run_helper sed -n '/  "Labels": {/, /  },/p' $PODMAN_TUI_LOG
    assert "$output" =~ "\"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\"" "expected \"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\" in pod inspect"
}

@test "pod remove" {
    # switch to pods view
    # select test pod from list
    # select remove command from pod commands dialog
    # confirm the operation on warnings dialog
    podman_tui_set_view "pods"
    podman_tui_select_item 0
    podman_tui_select_pod_cmd "remove"
    podman_tui_send_inputs "Enter"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --format "{{ .Name }}" --filter "name=${TEST_POD_NAME}$"
    assert "$output" == "" "expected $TEST_POD_NAME pod removal"
}

@test "pod prune" {
    podman pod create --name $TEST_POD_NAME --label $TEST_LABEL || echo done
    podman pod start $TEST_POD_NAME || echo done
    podman pod stop $TEST_POD_NAME || echo done
    sleep $TEST_TIMEOUT_LOW

    # switch to pods view
    # select prune command from pod commands dialog
    # confirm the operation on warnings dialog
    podman_tui_set_view "pods"
    podman_tui_select_pod_cmd "prune"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman pod ls --format "{{ .Name }}" --filter "name=${TEST_POD_NAME}$"
    assert "$output" == "" "expected at least $TEST_POD_NAME pod removal"
}
