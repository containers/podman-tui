# Podman TUI Documents

## Key Binding

podman-tui uses following keyboard keys for different actions:

| Action                           | Command |
| -------------------------------- | ------- |
| To view pods list page           | F1      |
| To view containers list page     | F2      |
| To view volumes list page        | F3      |
| To view images list page         | F4      |
| To view networks list page       | F5      |
| To view system page              | F6      |
| Lunch page command dialog        | Enter   |
| Close a dialog                   | Esc     |
| Switch between interface widgets | Tab     |

---

## Available Podman Commands

### pod

| COMMAND  | DESCRIPTION |
| -------- | ----------- |
| create   | create a new pod
| inpsect  | display information describing the seleceted pod
| kill     | send SIGTERM singnal to containers in the pod
| pause    | pause the selected pod
| prune    | remove all stopped pods and their containers
| restart  | restart the selected pod
| rm       | remove the selected pod
| start    | start  the selected pod
| stop     | stop th the selected pod
| top      | display the running processes of the pod's containers
| unpause  | unpause  the selected pod


### container

| COMMAND  | DESCRIPTION |
| -------- | ----------- |
| create   | create a new container but do not start
| diff     | inspect changes to the selected container's file systems
| inpsect  | display the configuration of a container
| kill     | kill the selected running container with a SIGKILL signal
| logs     | fetch the logs of the selected container
| pause    | pause all the processes in the selected container
| port     | list port mappings for the selected container
| prune    | remove all non running containers
| rename   | rename the selected container
| rm       | remove the selected container
| start    | start the selected containers
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