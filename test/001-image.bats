#!/usr/bin/env bats
#
# podman-tui images view functionality tests
#

load helpers
load helpers_tui


@test "image search and pull" {
    check_skip "image_search"

    podman image rm busybox || echo done

    # switch to images view
    # select pull/search command from images commands dialog
    # fillout search input field and press enter
    # wait for search operation
    # press KeyDown and then Enter to pull the image
    # wait for image pull operation
    podman_tui_set_view "images"
    podman_tui_select_image_cmd "pull"
    podman_tui_send_inputs "busybox" "Enter"
    sleep $TEST_TIMEOUT_HIGH
    podman_tui_send_inputs "Down" "Enter"
    sleep $TEST_TIMEOUT_HIGH

    run_helper podman image ls busybox --format "{{ .Repository }}"
    assert "$output" =~ "docker.io/library/busybox" "expected image"
}

@test "image save" {
    check_skip "image_save"

    podman image pull docker.io/library/busybox || echo done
    [ -f "${TEST_IMAGE_SAVE_PATH}" ] && /bin/rm -rf $TEST_IMAGE_SAVE_PATH

    # switch to images view
    # select busybox image
    # select save command from images commands dialog
    # fillout output path
    # select compressed options
    # go to save button and press enter

    podman_tui_set_view "images"
    podman_tui_select_item 0
    podman_tui_select_image_cmd "save"

    podman_tui_send_inputs $TEST_IMAGE_SAVE_PATH "Tab"
    podman_tui_send_inputs "Space" "Tab" "Tab" "Tab" "Tab"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_MEDIUM

    run_helper ls ${TEST_IMAGE_SAVE_PATH} 2> /dev/null
    assert "$output" == "$TEST_IMAGE_SAVE_PATH" "expected $TEST_IMAGE_SAVE_PATH exists"
}

@test "image import" {
    check_skip "image_import"

    /bin/rm -rf ${TEST_IMAGE_SAVE_PATH} || echo done
    podman image save -o ${TEST_IMAGE_SAVE_PATH} busybox:latest || echo done
    # switch to images view
    # select import command from images commands dialog
    # fillout import path field
    # fillout reference field
    # go to import button and press enter

    podman_tui_set_view "images"
    podman_tui_select_image_cmd "import"

    podman_tui_send_inputs $TEST_IMAGE_SAVE_PATH
    podman_tui_send_inputs "Tab" "Tab"
    podman_tui_send_inputs "${TEST_NAME}_image_imported"
    podman_tui_send_inputs "Tab"
    podman_tui_send_inputs "localhost/${TEST_NAME}_image_imported:latest"
    podman_tui_send_inputs "Tab" "Tab" "Enter"
    sleep $TEST_TIMEOUT_MEDIUM

    run_helper podman image ls ${TEST_NAME}_image_imported --format "{{ .Repository }}:{{ .Tag }}"
    assert "$output" =~ "localhost/${TEST_NAME}_image_imported" "expected image"
}

@test "image build" {
    check_skip "image_build"

    podman image pull docker.io/library/busybox || echo done
    podman image rm ${TEST_IMAGE_BUILD_REPOSITORY}/${TEST_IMAGE_BUILD_CONTEXT_DIR} || echo done

    # switch to images view
    # select build command from images commands dialog
    # fillout Context dir field
    # fillout image tag field
    # fillout image repository field
    # go to build button and press enter
    # wait for image build
    podman_tui_set_view "images"
    podman_tui_select_image_cmd "build"
    podman_tui_send_inputs ${TEST_IMAGE_BUILD_CONTEXT_DIR}
    podman_tui_send_inputs "Tab" "Tab" "Tab"
    podman_tui_send_inputs ${TEST_IMAGE_BUILD_TAG}
    podman_tui_send_inputs "Tab"
    podman_tui_send_inputs ${TEST_IMAGE_BUILD_REPOSITORY}
    podman_tui_send_inputs "Tab" "Tab"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_MEDIUM
    podman_tui_send_inputs "Tab" "Enter"

    run_helper podman image ls ${TEST_IMAGE_BUILD_TAG} --format "{{ .Repository }}:{{ .Tag }}"
    assert "$output" =~ "${TEST_IMAGE_BUILD_REPOSITORY}/${TEST_IMAGE_BUILD_TAG}" "expected image"
}

@test "image diff" {
    check_skip "image_diff"

    # switch to images view
    # select busybox image from list
    # select diff command from image commands dialog
    # close busybox image diff result message dialog
    podman_tui_set_view "images"
    podman_tui_select_item 2
    podman_tui_select_image_cmd "diff"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab" "Enter"

    run_helper grep -w 'A /var' $PODMAN_TUI_LOG
    assert "$output" =~ 'A /var' "expected 'A /var' in the logs"
}

@test "image history" {
    check_skip "image_history"

    # switch to images view
    # select busybox image from list
    # select history command from image commands dialog
    # close busybox image history result message dialog
    podman_tui_set_view "images"
    podman_tui_select_item 2
    podman_tui_select_image_cmd "history"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab" "Enter"

    run_helper egrep -w "\[\[$image_id.*BusyBox.*" $PODMAN_TUI_LOG
    assert "$output" =~ "$image_id" "expected $image_id in the history logs"
}

@test "image inspect" {
    check_skip "image_inspect"

    image_id=$(podman image ls --sort repository --noheading | nl -v 0 | grep 'busybox ' | awk '{print $4}')

    # switch to images view
    # select busybox image from list
    # select inspect command from image commands dialog
    # close busybox image inspect result message dialog
    podman_tui_set_view "images"
    podman_tui_select_item 2
    podman_tui_select_image_cmd "inspect"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Enter"

    run_helper sed -n '/  "RepoTags": \[/, /  \],/p' $PODMAN_TUI_LOG
    assert "$output" =~ "docker.io/library/busybox:latest" "expected RepoTag: [\"docker.io/library/busybox:latest\"] in the logs"
}

@test "image tag" {
    check_skip "image_tag"

    # switch to images view
    # select busybox image from list
    # select tag command from image commands dialog
    podman_tui_set_view "images"
    podman_tui_select_item 2
    podman_tui_select_image_cmd "tag"
    podman_tui_send_inputs "$TEST_IMAGE_TAG_NAME" "Tab" "Tab" "Enter"
    sleep $TEST_TIMEOUT_LOW

    run_helper podman image ls $TEST_IMAGE_TAG_NAME --format "{{ .Repository }}"
    assert "$output" =~ "$TEST_IMAGE_TAG_NAME" "expected tagged image $TEST_IMAGE_TAG_NAME"
}

@test "image untag" {
    check_skip "image_tag"

    # switch to images view
    # select busybox image from list
    # select untag command from image commands dialog
    # press "Tab" 2 times and "Enter" to untag busybox image
    podman_tui_set_view "images"
    podman_tui_select_item 2
    sleep $TEST_TIMEOUT_LOW
    podman_tui_select_image_cmd "untag"
    podman_tui_send_inputs "Tab" "Tab" "Enter"
    sleep $TEST_TIMEOUT_LOW

    untagged_umage=$(podman image ls --format '{{ .Repository }}')
    assert "$untagged_umage" !~ "$TEST_IMAGE_TAG_NAME" "expected $TEST_IMAGE_TAG_NAME not to be in the list"

}

@test "image remove" {
    check_skip "image_remove"

    run_helper podman image ls  --format "'{{ .Repository }}'"
    before="$output"

    # switch to images view
    # select <none> image from list
    # select remove command from image commands dialog
    # confirm image removal process and close warnings dialog
    # wait for image removal operation
    # close image removal result message dialog
    podman_tui_set_view "images"
    podman_tui_select_item 2
    podman_tui_select_image_cmd "remove"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW
    podman_tui_send_inputs "Tab" "Enter"

    # check if busybox image has been removed
    run_helper podman image ls --format "{{ .Repository }}"
    assert "$output" !~ "$before" "expected <none> image removal"
}

@test "image prune" {
    check_skip "image_prune"

    podman image pull busybox

    # switch to images view
    # select prune command from image commands dialog
    # confirm the operation on warnings dialog
    podman_tui_set_view "images"
    podman_tui_select_image_cmd "prune"
    podman_tui_send_inputs "Enter"
    sleep $TEST_TIMEOUT_LOW

    # check if busybox image has been removed
    run_helper podman image ls --format "{{ .Repository }}" --filter "reference=busybox"
    assert "$output" == "" "expected at least busybox image removal"
}
