# -*- bash -*-

TEST_IMAGE_TAG_NAME="podman_tui_busybox_tag01"
TEST_NAME="podman_tui_test"
TEST_VOLUME_NAME="${TEST_NAME}_vol01"
TEST_NETWORK_NAME="${TEST_NAME}_net01"
TEST_POD_NAME="${TEST_NAME}_pod01"
TEST_CONTAINER_NAME="${TEST_NAME}_container01"
TEST_CONTAINER_PORT="8888:80"
TEST_LABEL_NAME="test"
TEST_LABEL_VALUE="$TEST_NAME"
TEST_LABEL="${TEST_LABEL_NAME}=${TEST_LABEL_VALUE}"

################
#  podman_tui_set_view  # switches to different podman-tui views
################
function podman_tui_set_view() {
  case $1 in
  "pods")
    run_helper tmux send-keys -t $TMUX_SESSION F1;;
  "containers")
    run_helper tmux send-keys -t $TMUX_SESSION F2;;
  "volumes")
    run_helper tmux send-keys -t $TMUX_SESSION F3;;
  "images")
    run_helper tmux send-keys -t $TMUX_SESSION F4;;
  "networks")
    run_helper tmux send-keys -t $TMUX_SESSION F5;;
  "system")
    run_helper tmux send-keys -t $TMUX_SESSION F6;;
  esac   
}

################
#  podman_tui_select_item  # selects item from main view (pods, container, ...)
################
function podman_tui_select_item() {
  local menu_index=$1
  local current_index=0
  while [[ $current_index -lt $menu_index ]]
  do
    run_helper tmux send-keys -t $TMUX_SESSION Down
    let current_index=current_index+1
  done

}

################
#  podman_tui_send_inputs  # sends inputs to focus primitive of tui
################
function podman_tui_send_inputs() {
  for key in $@
  do
    run_helper tmux send-keys -t $TMUX_SESSION "$key"
  done
    
}

################
#  podman_tui_select_image_cmd # selects image command from cmd dialog
################
function podman_tui_select_image_cmd() {
  local menu_index=0
  case $1 in
  "diff")
    menu_index=0;;
  "history")
    menu_index=1;;
  "inspect")
    menu_index=2;;
  "prune")
    menu_index=3;;
  "remove")
    menu_index=4;;
  "pull")
    menu_index=5;;
  "tag")
    menu_index=6;;
  "untag")
    menu_index=7;;
  esac

  podman_tui_select_menu $menu_index
  
}

################
#  podman_tui_select_volume_cmd # selects volume command from cmd dialog
################
function podman_tui_select_volume_cmd() {
  local menu_index=0
  case $1 in
  "create")
    menu_index=0;;
  "inspect")
    menu_index=1;;
  "prune")
    menu_index=2;;
  "remove")
    menu_index=3;;
  esac

  podman_tui_select_menu $menu_index
}

################
#  podman_tui_select_network_cmd # selects network command from cmd dialog
################
function podman_tui_select_network_cmd() {
  local menu_index=0
  case $1 in
  "create")
    menu_index=0;;
  "inspect")
    menu_index=1;;
  "prune")
    menu_index=2;;
  "remove")
    menu_index=3;;
  esac

  podman_tui_select_menu $menu_index
}

################
#  podman_tui_select_pod_cmd # selects pod command from cmd dialog
################
function podman_tui_select_pod_cmd() {
  local menu_index=0

  case $1 in
  "create")
    menu_index=0;;
  "inspect")
    menu_index=1;;
  "kill")
    menu_index=2;;
  "pause")
    menu_index=3;;
  "prune")
    menu_index=4;;
  "restart")
    menu_index=5;;
  "remove")
    menu_index=6;;
  "start")
    menu_index=7;;
  "stop")
    menu_index=8;;
  "top")
    menu_index=9;;
  "unpause")
    menu_index=10;;
  esac

  podman_tui_select_menu $menu_index
}

################
#  podman_tui_select_container_cmd # selects container command from cmd dialog
################
function podman_tui_select_container_cmd() {
  local menu_index=0

  case $1 in
  "create")
    menu_index=0;;
  "diff")
    menu_index=1;;
  "exec")
    menu_index=2;;
  "inspect")
    menu_index=3;;
  "kill")
    menu_index=4;;
  "logs")
    menu_index=5;;
  "pause")
    menu_index=6;;
  "port")
    menu_index=7;;
  "prune")
    menu_index=8;;
  "rename")
    menu_index=9;;
  "remove")
    menu_index=10;;
  "start")
    menu_index=11;;
  "stop")
    menu_index=12;;
  "top")
    menu_index=13;;
  "unpause")
    menu_index=14;;
  esac

  podman_tui_select_menu $menu_index
}


################
#  podman_tui_select_menu # selects menu from menu dialog
################
function podman_tui_select_menu() {
  local menu_index=$1
  local current_index=0

  run_helper tmux send-keys -t $TMUX_SESSION Enter
  while [[ $current_index -lt $menu_index ]]
  do
    run_helper tmux send-keys -t $TMUX_SESSION Down
    let current_index=current_index+1
  done
  run_helper tmux send-keys -t $TMUX_SESSION Enter
}
