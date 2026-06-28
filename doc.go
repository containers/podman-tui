/*
podman-tui is a Terminal User Interface to interact with the podman (v6.x) and
its using podman API bindings to communicate with podman environment (unix socket).

Build and run from source:

  - clone repository
  - # make binary
  - # systemctl --user start podman.socket
  - # export TERM=xterm-256color
  - # ./bin/podman-tui

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
  - m: display command menu
  - s: display sort menu
  - l: switch to next screen
  - h: switch to previous screen
  - k: move up
  - j: move down
  - Ctrl-C: exit
  - Delete: deleted focus/selected item
  - Up/Down: move up/down
  - Left/Right: move left/right
  - Page Up/Down: scroll up/down
*/
package main
