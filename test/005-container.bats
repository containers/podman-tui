#!/usr/bin/env bats
#
# podman-tui containers view functionality tests
#

load helpers
load helpers_tui

@test "container run" {
    podman container rm -f $TEST_CONTAINER_NAME || echo done

    buysbox_image=$(podman image ls --sort repository --format "{{ .Repository }}" --filter "reference=docker.io/library/busybox")
    if [ "${buysbox_image}" == "" ] ; then
        podman image pull docker.io/library/busybox
    fi

    image_index=$(podman image ls --sort repository --noheading | nl -v 1 | grep 'busybox ' | awk '{print $1}')

    podman_tui_set_view "containers"
    podman_tui_select_container_cmd "run"
    podman_tui_send_inputs $TEST_CONTAINER_NAME "Tab" "$TEST_CONTAINER_RUN_CMD" "Tab"
    podman_tui_send_inputs "Down"
    podman_tui_select_item $image_index
    podman_tui_send_inputs "Enter"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Space"
    podman_tui_send_inputs "Tab" "Tab" "Space"
    podman_tui_send_inputs "Tab" "Space"
    podman_tui_send_inputs "Tab" "Space"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Enter"
    sleep $TEST_TIMEOUT_HIGH

    cnt_status=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .State.Status }}")
    assert "$cnt_status" =~ "running" "expected container status to match: running"

    podman container stop $TEST_CONTAINER_NAME

    run_helper podman container ls --all --filter "name=${TEST_CONTAINER_NAME}$" --noheading
    assert "$output" == "" "expected $TEST_CONTAINER_NAME to be removed"
}

@test "container create (privileged, timeout, remove)" {
    podman container rm -f $TEST_CONTAINER_NAME || echo done

    buysbox_image=$(podman image ls --sort repository --format "{{ .Repository }}" --filter "reference=docker.io/library/busybox")
    if [ "${buysbox_image}" == "" ] ; then
        podman image pull docker.io/library/busybox
    fi

    image_index=$(podman image ls --sort repository --noheading | nl -v 1 | grep 'busybox ' | awk '{print $1}')

    # switch to containers view
    # select create command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_container_cmd "create"

    # fillout name field
    # select image from dropdown widget
    # select privileged
    # set timeout to 10
    podman_tui_send_inputs $TEST_CONTAINER_NAME "Tab" "Tab"
    podman_tui_send_inputs "Down"
    podman_tui_select_item $image_index
    podman_tui_send_inputs "Enter" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Space" "Tab" "Space" "Tab" "$TEST_CONTAINER_TIMEOUT"
    podman_tui_send_inputs "Tab" "Tab" "Tab"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    cnt_status=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .State.Status }}")
    cnt_annotations=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .Config.Annotations }}")
    cnt_timeout=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .Config.Timeout }}")

    assert "$cnt_status" =~ "created" "expected container status to match: created"
    assert "$cnt_annotations" =~ '"io.podman.annotations.autoremove":"TRUE"' "expected container annotations to include: io.podman.annotations.autoremove:TRUE"
    assert "$cnt_annotations" =~ '"io.podman.annotations.autoremove":"TRUE"' "expected container annotations to include: io.podman.annotations.privileged:TRUE"
    assert "$cnt_timeout" =~ "$TEST_CONTAINER_TIMEOUT" "expected container config timeout to match: $TEST_CONTAINER_TIMEOUT"
}

@test "container create (environment page)" {
    podman container rm -f $TEST_CONTAINER_NAME || echo done

    buysbox_image=$(podman image ls --sort repository --format "{{ .Repository }}" --filter "reference=docker.io/library/busybox")
    if [ "${buysbox_image}" == "" ] ; then
        podman image pull docker.io/library/busybox
    fi

    image_index=$(podman image ls --sort repository --noheading | nl -v 1 | grep 'busybox ' | awk '{print $1}')

    # switch to containers view
    # select create command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_container_cmd "create"

    # fillout name field
    # select image from dropdown widget
    podman_tui_send_inputs $TEST_CONTAINER_NAME "Tab" "Tab"
    podman_tui_send_inputs "Down"
    podman_tui_select_item $image_index
    podman_tui_send_inputs "Enter" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab"
    sleep $TEST_TIMEOUT_LOW

    # switch to environmen page
    podman_tui_send_inputs "Down" "Tab"
    podman_tui_send_inputs "$TEST_CONTAINER_WORKDIR" "Tab"
    podman_tui_send_inputs "$TEST_CONTAINER_ENV1" "Space" "$TEST_CONTAINER_ENV2"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "$TEST_CONTAINER_UMASK"
    podman_tui_send_inputs "Tab" "Tab"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    cnt_workdir=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .Config.WorkingDir }}")
    cnt_vars=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .Config.Env }}")
    cnt_umask=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .Config.Umask }}")

    assert "$cnt_workdir" =~ "$TEST_CONTAINER_WORKDIR" "expected container work dir to match: $TEST_CONTAINER_WORKDIR"
    assert "$cnt_umask" =~ "$TEST_CONTAINER_UMASK" "expected container umask to match: $TEST_CONTAINER_UMASK"
    assert "$cnt_vars" =~ "$TEST_CONTAINER_ENV1" "expected container env to match: $TEST_CONTAINER_ENV1"
    assert "$cnt_vars" =~ "$TEST_CONTAINER_ENV2" "expected container env to match: $TEST_CONTAINER_ENV2"
}

@test "container create (resource page)" {
    podman container rm -f $TEST_CONTAINER_NAME || echo done

    buysbox_image=$(podman image ls --sort repository --format "{{ .Repository }}" --filter "reference=docker.io/library/busybox")
    if [ "${buysbox_image}" == "" ] ; then
        podman image pull docker.io/library/busybox
    fi

    image_index=$(podman image ls --sort repository --noheading | nl -v 1 | grep 'busybox ' | awk '{print $1}')

    # switch to containers view
    # select create command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_container_cmd "create"

    # fillout name field
    # select image from dropdown widget
    podman_tui_send_inputs $TEST_CONTAINER_NAME "Tab" "Tab"
    podman_tui_send_inputs "Down"
    podman_tui_select_item $image_index
    podman_tui_send_inputs "Enter" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab" "Tab"
    sleep $TEST_TIMEOUT_LOW

    # switch to environmen page
    podman_tui_send_inputs "Down" "Down" "Down" "Down" "Down" "Down" "Down" "Down" "Down" "Tab"
    podman_tui_send_inputs "$TEST_CONTAINER_MEMORY" "Tab" "$TEST_CONTAINER_MEMORY_RESERV" "Tab"
    podman_tui_send_inputs "$TEST_CONTAINER_MEMORY_SWAP" "Tab" "$TEST_CONTAINER_MEMORY_SWAPPINESS"
    podman_tui_send_inputs "Tab" "Tab" "$TEST_CONTAINER_CPU_SHARES"
    podman_tui_send_inputs "Tab" "$TEST_CONTAINER_CPU_PERIOD"
    podman_tui_send_inputs "Tab" "Tab" "$TEST_CONTAINER_CPU_QUOTA"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "$TEST_CONTAINER_SHM_SIZE" "Tab"
    podman_tui_send_inputs "$TEST_CONTAINER_SHM_SIZE_SYSTYEMD"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab" "Tab"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    cnt_memory=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .HostConfig.Memory }}")
    cnt_memory_reserv=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .HostConfig.MemoryReservation }}")
    cnt_memory_swap=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .HostConfig.MemorySwap }}")
    cnt_memory_swappiness=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .HostConfig.MemorySwappiness }}")
    cnt_cpu_shares=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .HostConfig.CpuShares }}")
    cnt_cpu_period=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .HostConfig.CpuPeriod }}")
    cnt_cpu_quota=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .HostConfig.CpuQuota }}")
    cnt_shm_size=$(podman container inspect $TEST_CONTAINER_NAME --format "{{ json .HostConfig.ShmSize }}")

    assert "$cnt_memory" =~ "$TEST_CONTAINER_MEMORY" "expected container memory to match: $TEST_CONTAINER_MEMORY"
    assert "$cnt_memory_reserv" =~ "$TEST_CONTAINER_MEMORY_RESERV" "expected container memory reservation to match: $TEST_CONTAINER_MEMORY_RESERV"
    assert "$cnt_memory_swap" =~ "$TEST_CONTAINER_MEMORY_SWAP" "expected container memory swap to match: $TEST_CONTAINER_MEMORY_SWAP"
    assert "$cnt_memory_swappiness" =~ "$TEST_CONTAINER_MEMORY_SWAPPINESS" "expected container memory swappiness to match: $TEST_CONTAINER_MEMORY_SWAPPINESS"
    assert "$cnt_cpu_shares" =~ "$TEST_CONTAINER_CPU_SHARES" "expected container cpu shares to match: $TEST_CONTAINER_CPU_SHARES"
    assert "$cnt_cpu_period" =~ "$TEST_CONTAINER_CPU_PERIOD" "expected container cpu period to match: $TEST_CONTAINER_CPU_PERIOD"
    assert "$cnt_cpu_quota" =~ "$TEST_CONTAINER_CPU_QUOTA" "expected container cpu quota to match: $TEST_CONTAINER_CPU_QUOTA"
    assert "$cnt_shm_size" =~ "$TEST_CONTAINER_SHM_SIZE" "expected container shm size to match: $TEST_CONTAINER_SHM_SIZE"
}

@test "container create (pod, network, volume, security options, health)" {
    httpd_image=$(podman image ls --sort repository --format "{{ .Repository }}" --filter "reference=docker.io/library/httpd")
    if [ "${httpd_image}" == "" ] ; then
        podman image pull docker.io/library/httpd
    fi

    podman pod rm -f $TEST_CONTAINER_POD_NAME || echo done
    podman container rm -f $TEST_CONTAINER_NAME || echo done
    podman container rm -f ${TEST_CONTAINER_NAME}_renamed || echo done
    podman network rm $TEST_CONTAINER_NETWORK_NAME || echo done
    podman volume rm $TEST_CONTAINER_VOLUME_NAME || echo done

    [ ! -d "${TEST_CONTAINER_MOUNT_SOURCE}" ] && mkdir $TEST_CONTAINER_MOUNT_SOURCE

    podman network create $TEST_CONTAINER_NETWORK_NAME || echo done
    podman volume create $TEST_CONTAINER_VOLUME_NAME || echo done
    podman pod create --name $TEST_CONTAINER_POD_NAME --network $TEST_CONTAINER_NETWORK_NAME --publish $TEST_CONTAINER_PORT || echo done
    sleep $TEST_TIMEOUT_LOW

    # get required pod, image, network and volume index for number of KeyDown stroke
    pod_index=$(podman pod ls --sort name --format "{{ .Name }}" | nl -v 1 | grep "$TEST_CONTAINER_POD_NAME" | awk '{print $1}')
    image_index=$(podman image ls --sort repository --noheading | nl -v 1 | grep 'httpd ' | awk '{print $1}')
    net_index=$(podman network ls -q | nl -v 1 | grep "$TEST_CONTAINER_NETWORK_NAME" | awk '{print $1}')


    # switch to containers view
    # select create command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_container_cmd "create"

    # fillout name field
    # select image from dropdown widget
    # select pod from dropdown widget
    # fillout label field
    podman_tui_send_inputs $TEST_CONTAINER_NAME "Tab" "Tab"
    podman_tui_send_inputs "Down"
    podman_tui_select_item $image_index
    podman_tui_send_inputs "Enter" "Tab"
    podman_tui_send_inputs "Down"
    podman_tui_select_item $pod_index
    podman_tui_send_inputs "Enter" "Tab"
    podman_tui_send_inputs $TEST_LABEL "Tab" "Tab" "Tab" "Tab" "Tab" "Tab"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab"
    sleep $TEST_TIMEOUT_LOW

    # switch to "health check"  create view
    podman_tui_send_inputs "Down" "Down" "Down" "Down" "Tab"
    podman_tui_send_inputs $TEST_CONTAINER_HEALTH_CMD "Tab" "Tab"
    podman_tui_send_inputs "Enter" "Down" "Down" "Enter"
    podman_tui_send_inputs "Tab" "Tab" "Tab"
    podman_tui_send_inputs $TEST_CONTAINER_HEALTH_INTERVAL
    podman_tui_send_inputs "Tab" "Tab"
    podman_tui_send_inputs $TEST_CONTAINER_HEALTH_RETRIES
    podman_tui_send_inputs "Tab" "Tab"
    podman_tui_send_inputs $TEST_CONTAINER_HEALTH_TIMEOUT
    podman_tui_send_inputs "Tab" "Tab" "Tab"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab"
    sleep $TEST_TIMEOUT_LOW

    # switch to "security options" create view
    podman_tui_send_inputs "Down" "Down" "Down" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Space" "Tab" "Tab" "Tab"

    # switch to "volumes settings" create view
    # select volume from dropdown widget
    podman_tui_send_inputs "Down" "Tab"
    podman_tui_send_inputs "${TEST_CONTAINER_VOLUME_NAME}:${TEST_CONTAINER_VOLUME_MOUNT_POINT}:rw"
    podman_tui_send_inputs "Tab" "Tab"
    podman_tui_send_inputs "type=bind,src=${TEST_CONTAINER_MOUNT_SOURCE},dst=${TEST_CONTAINER_MOUNT_DEST}"
    sleep $TEST_TIMEOUT_LOW

    # go to "Create" button and press Enter
    podman_tui_send_inputs "Tab" "Tab" "Enter"
    sleep $TEST_TIMEOUT_LOW

    # get created container information
    container_information=$(podman container ls --all --pod --filter "name=${TEST_CONTAINER_NAME}$" --format \
    "{{ json .PodName }}|{{ json .Networks }}|{{ json .Mounts }}|{{ json .Image }}|{{ json .Ports }}|{{ json .Labels }}")

    cnt_pod_name=$(echo $container_information | awk -F '|' '{print $1}')
    cnt_networks=$(echo $container_information | awk -F '|' '{print $2}')
    cnt_mounts=$(echo $container_information | awk -F '|' '{print $3}')
    cnt_image_name=$(echo $container_information | awk -F '|' '{print $4}')
    cnt_ports=$(echo $container_information | awk -F '|' '{print $5}')
    cnt_labels=$(echo $container_information | awk -F '|' '{print $6}')

    host_port=$(echo $TEST_CONTAINER_PORT | awk -F: '{print $1}')
    cnt_port=$(echo $TEST_CONTAINER_PORT | awk -F: '{print $2}')
    cnt_port_str="$host_port->$cnt_port/tcp"

    cnt_security_opt=$(podman container inspect ${TEST_CONTAINER_NAME} --format "{{ json .HostConfig.SecurityOpt }}")

    cnt_healthcheck_interval=$(podman container inspect ${TEST_CONTAINER_NAME} --format "{{ json .Config.Healthcheck.Interval }}")
    cnt_healthcheck_timeout=$(podman container inspect ${TEST_CONTAINER_NAME} --format "{{ json .Config.Healthcheck.Timeout }}")
    cnt_healthcheck_retires=$(podman container inspect ${TEST_CONTAINER_NAME} --format "{{ json .Config.Healthcheck.Retries }}")
    cnt_hltcheck_interval="$(echo $TEST_CONTAINER_HEALTH_INTERVAL | sed 's/s//')000000000"
    cnt_gltcheck_timeout="$(echo $TEST_CONTAINER_HEALTH_TIMEOUT | sed 's/s//')000000000"

    assert "$cnt_pod_name" =~ "$TEST_CONTAINER_POD_NAME" "expected container pod: $TEST_CONTAINER_POD_NAME"

    assert "$cnt_mounts" =~ "$TEST_CONTAINER_VOLUME_MOUNT_POINT" "expected container volume mount point: $TEST_CONTAINER_VOLUME_MOUNT_POINT"
    assert "$cnt_mounts" =~ "$TEST_CONTAINER_MOUNT_DEST" "expected container mount point: $TEST_CONTAINER_MOUNT_DEST"

    assert "$cnt_image_name" =~ "$httpd_image" "expected container image name: $httpd_image"
    assert "$cnt_ports" =~ "$cnt_port_str" "expected container port: $cnt_port_str"
    assert "$cnt_security_opt" =~ "no-new-privileges" "expected no-new-privileges in container security options"

    # heathcheck tests
    assert "$cnt_healthcheck_interval" =~ "$cnt_hltcheck_interval" "expected healthcheck interval to mach $cnt_hltcheck_interval"
    assert "$cnt_healthcheck_timeout" =~ "$cnt_gltcheck_timeout" "expected healthcheck timeout to mach $cnt_gltcheck_timeout"
    assert "$cnt_healthcheck_retires" =~ "$TEST_CONTAINER_HEALTH_RETRIES" "expected healthcheck retries to mach $TEST_CONTAINER_HEALTH_RETRIES"

}

@test "container commit" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select container from the list
    # select commit command from container commands dialog
    # fillout image input field
    # go to commit button and press enter
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "commit"

    podman_tui_send_inputs $TEST_CONTAINER_COMMIT_IMAGE_NAME
    podman_tui_send_inputs Tab Tab Tab Tab
    podman_tui_send_inputs Tab Tab Tab Tab
    podman_tui_send_inputs Enter
    sleep $TEST_TIMEOUT_HIGH
    run_helper podman image ls ${TEST_CONTAINER_COMMIT_IMAGE_NAME} --format "{{ .Repository }}"
    assert "$output" =~ "localhost/${TEST_CONTAINER_COMMIT_IMAGE_NAME}" "expected image"
}

@test "container start" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select start command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "start"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman container ls --all --filter="name=${TEST_CONTAINER_NAME}$" --format "{{ .Status }}"
    assert "$output" =~ "Up" "expected $TEST_CONTAINER_NAME to be up"

}

@test "container checkpoint" {
    podman container create --name ${TEST_CONTAINER_CHECKPOINT_NAME} docker.io/library/httpd
    podman container start ${TEST_CONTAINER_CHECKPOINT_NAME}

    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_CHECKPOINT_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select checkpoint command from container commands dialog
    # fillout information
    # go to checkpoint button and Enter

    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "checkpoint"

    podman_tui_send_inputs "Tab"
    podman_tui_send_inputs "~/${TEST_CONTAINER_CHECKPOINT_NAME}_dump.tar"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Enter"

    sleep $TEST_TIMEOUT_HIGH

    run_helper ls ~/${TEST_CONTAINER_CHECKPOINT_NAME}_dump.tar 2>/dev/null || echo -e '\c'
    assert "$output" == "/root/${TEST_CONTAINER_CHECKPOINT_NAME}_dump.tar" "expected tar file to be created"

}

@test "containre restore" {
    # switch to containers view
    # selec restore command from container commands dialog
    # filleout information
    # go to restore button and Enter

    podman_tui_set_view "containers"
    podman_tui_select_container_cmd "restore"
    podman_tui_send_inputs "Tab" "Tab"
    podman_tui_send_inputs ${TEST_CONTAINER_CHECKPOINT_NAME}_restore
    podman_tui_send_inputs "Tab" "Tab"
    podman_tui_send_inputs "~/${TEST_CONTAINER_CHECKPOINT_NAME}_dump.tar"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Tab" "Tab" "Enter"

    sleep $TEST_TIMEOUT_HIGH
    run_helper podman container ls --all --format "{{ .Names }}"
    assert "$output" =~ "${TEST_CONTAINER_CHECKPOINT_NAME}_restore" "expected container to be restored"
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
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "echo Space test Space > Space a.txt" "Enter"
    podman_tui_send_inputs "Tab" "Enter"
    sleep $TEST_TIMEOUT_LOW

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
    sleep $TEST_TIMEOUT_LOW

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
    sleep $TEST_TIMEOUT_MEDIUM

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
    sleep $TEST_TIMEOUT_LOW

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
    sleep $TEST_TIMEOUT_LOW

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
    sleep $TEST_TIMEOUT_LOW

    run_helper podman container ls --all --filter="name=${TEST_CONTAINER_NAME}$" --format "{{ .Status }}"
    assert "$output" =~ "Paused" "expected $TEST_CONTAINER_NAME to be paused"
}

@test "container unpause" {
    container_index=$(podman container ls --all --format "{{ .Names }}" | sort | nl -v 0 | grep "$TEST_CONTAINER_NAME" | awk '{print $1}')

    # switch to containers view
    # select test container from list
    # select unpause command from container commands dialog
    podman_tui_set_view "containers"
    podman_tui_select_item $container_index
    podman_tui_select_container_cmd "unpause"
    sleep $TEST_TIMEOUT_LOW

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
    sleep $TEST_TIMEOUT_LOW

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
    sleep $TEST_TIMEOUT_LOW

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
    sleep $TEST_TIMEOUT_LOW

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
    sleep $TEST_TIMEOUT_LOW

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
    sleep $TEST_TIMEOUT_MEDIUM

    run_helper podman container ls --all --filter "name=${TEST_CONTAINER_NAME}$" --noheading
    assert "$output" == "" "expected $TEST_CONTAINER_NAME to be removed"
}
