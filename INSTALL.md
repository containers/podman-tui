## Building From Source

podman-tui is using go version >= 1.17. 
 1. Clone the repo
 2. Install [dependencies](./CONTRIBUTING.md/#prerequisite-before-build)
 3. Build

      ```shell
      make binary
      ```
 4. Run podman-tui

      ```shell
      ./bin/podman-tui
      ```

Run `sudo make install` if you want to install the binary on the node.

---

## Installing from the AUR (Arch Linux)

There are 2 packages present in the Arch User Repository:
- `podman-tui`
- `podman-tui-git`

Either one can be installed by using your prefered AUR helper, such as:

```shell
yay -S podman-tui
```

or manually:

```shell
git clone https://aur.archlinux.org/podman-tui.git
cd podman-tui
makepkg -si
```

