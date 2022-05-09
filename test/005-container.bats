#!/usr/bin/env bats
#
# podman-tui containers view functionality tests
#

load helpers
load helpers_tui

@test "container create" {
    podman pod rm -f $TEST_CONTAINER_POD_NAME || echo done
    podman container rm -f $TEST_CONTAINER_NAME || echo done
    podman container rm -f ${TEST_CONTAINER_NAME}_renamed || echo done
    podman network rm $TEST_CONTAINER_NETWORK_NAME || echo done
    podman volume rm $TEST_CONTAINER_VOLUME_NAME || echo done
    
    httpd_image=$(podman image ls --sort repository --format "{{ .Repository }}" --filter "reference=docker.io/library/httpd")
    if [ "${httpd_image}" == "" ] ; then 
        podman image pull docker.io/library/httpd
    fi

    podman network create $TEST_CONTAINER_NETWORK_NAME || echo done
    podman volume create $TEST_CONTAINER_VOLUME_NAME || echo done
    podman pod create --name $TEST_CONTAINER_POD_NAME || echo done
    sleep 2

    # get required pod, image, network and volume index for number of KeyDown stroke
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 1 | grep "$TEST_CONTAINER_POD_NAME" | awk '{print $1}')
    image_index=$(podman image ls --sort repository --noheading | nl -v 1 | grep 'httpd ' | awk '{print $1}')
    net_index=$(podman network ls -q | nl -v 1 | grep "$TEST_CONTAINER_NETWORK_NAME" | awk '{print $1}')
    vol_index=$(podman volume ls -q | nl -v 1 | grep "$TEST_CONTAINER_VOLUME_NAME" | awk '{print $1}')


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
    podman_tui_send_inputs $TEST_LABEL "Tab" "Tab" "Tab" "Tab"
    sleep 1

    # switch to "network settings" create view
    # select network from dropdown widget
    podman_tui_send_inputs "Down" "Down" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Down"
    podman_tui_select_item $net_index
    podman_tui_send_inputs "Enter"
    podman_tui_send_inputs "Tab" "Tab"
    sleep 1
    
    # switch to "ports settings" create view
    # fillout "publish ports" field
    podman_tui_send_inputs "Tab" "Down" "Tab"
    podman_tui_send_inputs $TEST_CONTAINER_PORT "Tab" "Tab" "Tab" "Tab" "Tab"
    sleep 1

    # switch to "security options" create view
    podman_tui_send_inputs "Down" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Space" "Tab" "Tab" "Tab"

    # switch to "volumes settings" create view
    # select volume from dropdown widget
    podman_tui_send_inputs "Down" "Tab" "Down"
    podman_tui_select_item $vol_index
    podman_tui_send_inputs "Enter"
    sleep 1

    # go to "Create" button and press Enter
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Enter"
    sleep 2

    # get created container information
    container_information=$(podman container ls --all --pod --filter "name=${TEST_CONTAINER_NAME}$" --format \
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
    
    cnt_security_opt=$(podman container inspect ${TEST_CONTAINER_NAME} --format "{{ json .HostConfig.SecurityOpt }}")

    assert "$cnt_pod_name" =~ "$TEST_CONTAINER_POD_NAME" "expected container pod: $TEST_CONTAINER_POD_NAME"
    assert "$cnt_mounts" =~ "$TEST_CONTAINER_VOLUME_NAME" "expected container volume: $TEST_CONTAINER_VOLUME_NAME"
    assert "$cnt_image_name" =~ "$httpd_image" "expected container image name: $httpd_image"
    assert "$cnt_ports" =~ "$cnt_port_str" "expected container port: $cnt_port_str"
    assert "$cnt_security_opt" =~ "no-new-privileges" "expected no-new-privileges in container security options"

}

@test "container start" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select start command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "start"
    sleep 2

    run_helper podman container ls --all --filter="name=${TEST_CONTAINER_NAME}$" --format "{{ .Status }}"
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
    podman_tui_send_inputs "Tab" "Space" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Enter"
    sleep 1
    podman_tui_send_inputs "echo Space test Space > Space a.txt" "Enter"
    podman_tui_send_inputs "Tab" "Enter"
    sleep 2

    run_helper podman container exec $TEST_CONTAINER_NAME cat a.txt

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
    sleep 2

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

    run_helper grep -w "/etc" $PODMAN_TUI_LOG
    assert "$output" =~ '/etc' "expected '/etc' in the logs"
}

@test "container top" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select top command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "top"
    sleep 2

    run_helper grep -w "USER PID PPID" $PODMAN_TUI_LOG
    assert "$output" =~ 'USER PID PPID' "expected 'USER PID PPID' in the logs"
}

@test "container port" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select port command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "port"
    sleep 2

    container_ports=$(podman container ls --all --filter="name=${TEST_CONTAINER_NAME}$" --format "{{ .Ports }}")
    run_helper grep -w "$container_ports" $PODMAN_TUI_LOG
    assert "$output" =~ "$container_ports" "expected container ports ($container_ports) in the log"
}

@test "container pause" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select pause command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "pause"
    sleep 2

    run_helper podman container ls --all --filter="name=${TEST_CONTAINER_NAME}$" --format "{{ .Status }}"
    assert "$output" == "Paused" "expected $TEST_CONTAINER_NAME to be paused"
}

@test "container unpause" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select unpause command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "unpause"
    sleep 2

    run_helper podman container ls --all --filter="name=${TEST_CONTAINER_NAME}$" --format "{{ .Status }}"
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
    sleep 2

    run_helper podman container ls --all --filter="name=${TEST_CONTAINER_NAME}$" --format "{{ .Status }}"
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
    sleep 2

    run_helper podman container ls --all --filter="name=${TEST_CONTAINER_NAME}$" --format "{{ .Status }}"
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
    sleep 2

    run_helper podman container ls --all --filter "name=${TEST_CONTAINER_NAME}$" --noheading
    assert "$output" == "" "expected $TEST_CONTAINER_NAME to be removed"
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
    sleep 2

    run_helper podman container ls --all --filter "name=${TEST_CONTAINER_NAME}_renamed$" --format "{{ .Names }}"
    assert "$output" == "${TEST_CONTAINER_NAME}_renamed" "expected ${TEST_CONTAINER_NAME}_renamed to be in the list"
}

@test "container prune" {
    podman container create --name $TEST_CONTAINER_NAME docker.io/library/httpd || echo done
    podman container start $TEST_CONTAINER_NAME || echo done
    podman container stop $TEST_CONTAINER_NAME || echo done

    # switch to containers view
    # select test container from list
    # select prune command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "prune"
    podman_tui_send_inputs "Enter"
    sleep 3

    run_helper podman container ls --all --filter "name=${TEST_CONTAINER_NAME}$" --noheading
    assert "$output" == "" "expected $TEST_CONTAINER_NAME to be removed"
}
