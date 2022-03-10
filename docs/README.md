# Podman TUI Documents

## Key Binding

podman-tui uses following keyboard keys for different actions:

| Action                           | Command   |
| -------------------------------- | --------- |
| Display command menu             | m         |
| Switch to next screen            | l         |
| Switch to previous screen        | h         |
| Move up                          | k         |
| Move down                        | j         |
| Exit application                 | Ctrl+c    |
| Close the active dialog          | Esc       |
| Switch between interface widgets | Tab       |
| Move up/down                     | Up/Down   |
| Scroll Up                        | Page Up   |
| Scroll Down                      | Page Down |
| Display help screen              | F1        |
| Display pods screen              | F2        |
| Display containers screen        | F3        |
| Display volumes screen           | F4        |
| Display images screen            | F5        |
| Display networks screen          | F6        |
| Display system screen            | F7        |

---

## Available Podman Commands

### pod

| COMMAND  | DESCRIPTION |
| -------- | ----------- |
| create   | create a new pod
| inpsect  | display information describing the seleceted pod
| kill     | send SIGTERM signal to containers in the pod
| pause    | pause the selected pod
| prune    | remove all stopped pods and their containers
| restart  | restart the selected pod
| rm       | remove the selected pod
| start    | start the selected pod
| stats    | display the live stream of resource usage
| stop     | stop th the selected pod
| top      | display the running processes of the pod's containers
| unpause  | unpause  the selected pod


### container

| COMMAND  | DESCRIPTION |
| -------- | ----------- |
| create   | create a new container but do not start
| diff     | inspect changes to the selected container's file systems
| exec     | execute the specified command inside a running container
| inpsect  | display the configuration of a container
| kill     | kill the selected running container with a SIGKILL signal
| logs     | fetch the logs of the selected container
| pause    | pause all the processes in the selected container
| port     | list port mappings for the selected container
| prune    | remove all non running containers
| rename   | rename the selected container
| rm       | remove the selected container
| start    | start the selected containers
| stats    | display the live stream of resource usage
| stop     | stop the selected containers
| top      | display the running processes of the selected container
| unpause  | unpause the selected container that was paused before

### volume

| COMMAND  | DESCRIPTION |
| -------- | ----------- |
| create   | create a new volume
| inspect  | display detailed volume's information
| prune    | remove all unused volumes
| rm       | remove the selected volume

### image

| COMMAND      | DESCRIPTION |
| ------------ | ----------- |
| diff         | inspect changes to the image's file systems
| history      | show history of the selected image
| inspect      | display the configuration of the selected image
| prune        | remove all unused images
| rm           | removes the selected  image from local storage
| search/pull  | search and pull image from registry
| tag          | add an additional name to the selected  image
| untag        | remove a name from the selected image


### network

| COMMAND  | DESCRIPTION |
| -------- | ----------- |
| create   | create a Podman CNI network
| inspect  | displays the raw CNI network configuration
| prune    | remove all unused networks
| rm       | remove a CNI networks

### system

| COMMAND    | DESCRIPTION |
| ---------- | ----------- |
| disk usage | display Podman related system information
| info       | display system information
| prune      | remove all unused pod, container, image and volume data