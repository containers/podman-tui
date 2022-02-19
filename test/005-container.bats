#!/usr/bin/env bats
#
# podman-tui containers view functionality tests
#

load helpers
load helpers_tui

@test "container create" {
    podman container rm -f $TEST_CONTAINER_NAME || echo done 
    
    httpd_image=$(podman image ls --sort repository --format "{{ .Repository }}" --filter "reference=docker.io/library/httpd")
    if [ "${httpd_image}" == "" ] ; then 
        podman image pull docker.io/library/httpd
    fi
    podman network create $TEST_NETWORK_NAME || echo done
    podman volume create $TEST_VOLUME_NAME || echo done
    podman pod create --name $TEST_POD_NAME || echo done

    # get required pod, image, network and volume index for number of KeyDown stroke
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 1 | grep "$TEST_POD_NAME" | awk '{print $1}')
    image_index=$(podman image ls --sort repository --noheading | nl -v 1 | grep 'httpd ' | awk '{print $1}')
    net_index=$(podman network ls -q | nl -v 1 | grep "$TEST_NETWORK_NAME" | awk '{print $1}')
    vol_index=$(podman volume ls -q | nl -v 1 | grep "$TEST_VOLUME_NAME" | awk '{print $1}')


    # switch to containers view
    # select create command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_container_cmd "create"
    
    # fillout name field
    # select image from dropdown widget
    # select pod from dropdown widget
    # fillout label field
    podman_tui_send_inputs $TEST_CONTAINER_NAME "Tab"
    podman_tui_send_inputs "Down" 
    podman_tui_select_item $image_index
    podman_tui_send_inputs "Enter" "Tab"
    podman_tui_send_inputs "Down"
    podman_tui_select_item $pod_index
    podman_tui_send_inputs "Enter" "Tab"

    podman_tui_send_inputs $TEST_LABEL "Tab" "Tab" "Tab"
    sleep 1

    # switch to "network settings" create view
    # select network from dropdown widget
    podman_tui_send_inputs "Tab" "Down" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Down"
    podman_tui_select_item $net_index
    podman_tui_send_inputs "Enter"
    podman_tui_send_inputs "Tab" "Tab"
    sleep 1
    
    # switch to "ports settings" create view
    # fillout "publish ports" field
    podman_tui_send_inputs "Tab" "Down" "Tab"
    podman_tui_send_inputs $TEST_CONTAINER_PORT "Tab" "Tab" "Tab" "Tab"
    sleep 1

    # switch to "volumes settings" create view
    # select volume from dropdown widget
    podman_tui_send_inputs "Tab" "Down" "Down" "Tab" "Down"
    podman_tui_select_item $vol_index
    podman_tui_send_inputs "Enter" "Down" "Enter"
    sleep 1

    # go to "Create" button and press Enter
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Enter"
    sleep 1

    # get created container information
    container_information=$(podman container ls --all --pod --filter "name=podman_tui_test_container01" --format \
    "{{ json .PodName }},{{ json .Networks }},{{ json .Mounts }},{{ json .Image }},{{ json .Ports }},{{ json .Labels }}")

    cnt_pod_name=$(echo $container_information | awk -F, '{print $1}')
    cnt_networks=$(echo $container_information | awk -F, '{print $2}')
    cnt_mounts=$(echo $container_information | awk -F, '{print $3}')
    cnt_image_name=$(echo $container_information | awk -F, '{print $4}')
    cnt_ports=$(echo $container_information | awk -F, '{print $5}')
    cnt_labels=$(echo $container_information | awk -F, '{print $6}')

    host_port=$(echo $TEST_CONTAINER_PORT | awk -F: '{print $1}')
    cnt_port=$(echo $TEST_CONTAINER_PORT | awk -F: '{print $2}')
    cnt_port_str="$host_port->$cnt_port/tcp"

    assert "$cnt_pod_name" "=~" "$TEST_POD_NAME" "expected container pod: $TEST_POD_NAME"
    #assert "$cnt_networks" "=~" "$TEST_NETWORK_NAME" "expected container network: $TEST_NETWORK_NAME"
    assert "$cnt_mounts" "=~" "$TEST_VOLUME_NAME" "expected container volume: $TEST_VOLUME_NAME"
    assert "$cnt_image_name" "=~" "$httpd_image" "expected container image name: $httpd_image"
    assert "$cnt_ports" "=~" "$cnt_port_str" "expected container port: $cnt_port_str"
    #assert "$cnt_labels" "=~" "\"$TEST_LABEL_NAME\":\"$TEST_LABEL_VALUE\"" "expected container port: $TEST_CONTAINER_PORT"

}

@test "container start" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select start command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "start"
    sleep 1

    run_podman container ls --all --filter="name=$TEST_CONTAINER_NAME" --format "'{{ .Status }}'"
    assert "$output" =~ "Up" "expected $TEST_CONTAINER_NAME to be up"

}

@test "container exec" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select exec command from container commands dialog
    # fillout command field
    # check interactive checkbox
    # check tty checkbox
    # go to Execute button and Enter
    # type "echo test > a.txt"
    # go to Close button and Enter
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "exec"
    podman_tui_send_inputs "/bin/bash"
    podman_tui_send_inputs "Tab" "Space" "Tab" "Space"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Enter"
    sleep 1
    podman_tui_send_inputs "echo Space test Space > Space a.txt" "Enter"
    podman_tui_send_inputs "Tab" "Enter"
    sleep 1

    run_podman container exec $TEST_CONTAINER_NAME cat a.txt

    assert "$output" =~ "test" "expected container a.txt file container test keyword"
}

@test "container inspect" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select inspect command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "inspect"
    sleep 1

    run_helper sed -n '/  "Labels": {/, /  },/p' $PODMAN_TUI_LOG

    assert "$output" =~ "\"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\"" "expected \"$TEST_LABEL_NAME\": \"$TEST_LABEL_VALUE\" to be in labels"
}

@test "container diff" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select diff command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "diff"
    sleep 6

    run_helper grep -w "/home" $PODMAN_TUI_LOG
    assert "$output" "=~" '/home' "expected '/home' in the logs"
}

@test "container top" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select top command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "top"
    sleep 1

    run_helper grep -w "USER PID PPID" $PODMAN_TUI_LOG
    assert "$output" "=~" 'USER PID PPID' "expected 'USER PID PPID' in the logs"
}

@test "container port" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select port command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "port"
    sleep 1

    container_ports=$(podman container ls --all --filter="name=$TEST_CONTAINER_NAME" --format "{{ .Ports }}")
    run_helper grep -w "$container_ports" $PODMAN_TUI_LOG
    assert "$output" "=~" "$container_ports" "expected container ports ($container_ports) in the log"
}

@test "container pause" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select pause command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "pause"
    sleep 1

    run_podman container ls --all --filter="name=$TEST_CONTAINER_NAME" --format "'{{ .Status }}'"
    assert "$output" =~ "paused" "expected $TEST_CONTAINER_NAME to be paused"
}

@test "container unpause" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select unpause command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "unpause"
    sleep 1

    run_podman container ls --all --filter="name=$TEST_CONTAINER_NAME" --format "'{{ .Status }}'"
    assert "$output" =~ "Up" "expected $TEST_CONTAINER_NAME to be Up"
}

@test "container stop" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select stop command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "stop"
    sleep 1

    run_podman container ls --all --filter="name=$TEST_CONTAINER_NAME" --format "'{{ .Status }}'"
    assert "$output" =~ "Exited" "expected $TEST_CONTAINER_NAME to be Up"
}

@test "container kill" {
    podman container start $TEST_CONTAINER_NAME || echo done
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select kill command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "kill"
    sleep 1

    run_podman container ls --all --filter="name=$TEST_CONTAINER_NAME" --format "'{{ .Status }}'"
    assert "$output" =~ "Exited" "expected $TEST_CONTAINER_NAME to be killed"
}

@test "container remove" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select remove command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "remove"
    podman_tui_send_inputs "Enter"
    sleep 1

    run_podman container ls --all --format "'{{ .Names }}'"
    assert "$output" !~ "$TEST_CONTAINER_NAME" "expected $TEST_CONTAINER_NAME to be removed"
}

@test "container rename" {

    podman container create --name $TEST_CONTAINER_NAME httpd
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select rename command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "rename"
    podman_tui_send_inputs ${TEST_CONTAINER_NAME}_renamed
    podman_tui_send_inputs "Tab" "Tab" "Enter"
    sleep 1
    
    run_podman container ls --all --filter "name=${TEST_CONTAINER_NAME}\$" --format "'{{ json .Names }}'"
    assert "$output" "!~" "${TEST_CONTAINER_NAME}" "expected ${TEST_CONTAINER_NAME} to be not in the list"

    run_podman container ls --all --filter "name=${TEST_CONTAINER_NAME}_renamed\$" --format "'{{ json .Names }}'"
    assert "$output" "=~" "${TEST_CONTAINER_NAME}_renamed" "expected ${TEST_CONTAINER_NAME}_renamed to be in the list"
}

@test "container prune" {
    podman container create --name $TEST_CONTAINER_NAME docker.io/library/httpd || echo done

    # switch to containers view
    # select test container from list
    # select prune command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "prune"
    podman_tui_send_inputs "Enter"
    sleep 1

    run_podman container ls --all --format "'{{ .Names }}'"
    assert "$output" !~ "$TEST_CONTAINER_NAME" "expected $TEST_CONTAINER_NAME to be removed"
}
