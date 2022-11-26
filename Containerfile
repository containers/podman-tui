FROM docker.io/library/busybox:latest

COPY ./bin/podman-tui /bin/podman-tui

VOLUME /ssh_keys/

ENV TERM=xterm-256color

ENTRYPOINT [ "/bin/podman-tui" ]
