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

## Installing Packaged Versions

### Arch Linux (AUR)

```shell
yay -S podman-tui
```
