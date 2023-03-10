/*
podman-tui is a Terminal User Interface to interact with the podman (v3.x) and
its using podman API bindings to communicate with podman environment (unix socket).

Build from source:

  - clone repository
  - # make build
  - # ./bin/podman-tui

Pre-run checks:

  - # systemctl --user start podman.socket
  - # export TERM=xterm-256color

Key bindings:

  - F1: view pods list page
  - F2: view container list page
  - F3: view volumes list page
  - F4: view images list page
  - F5: view networks list page
  - F6: view system page
  - Enter: lunch page command dialog
  - Esc: close a dialog
  - Tab: switch between interface widgets
*/
package main
