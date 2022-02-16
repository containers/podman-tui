VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
    config.vm.hostname = "fedora34"
    config.vm.box = "fedora/34-cloud-base"
    config.vm.box_version = "34.20210423.0"
    config.vm.provision "shell", inline: "mkdir -p /home/vagrant/go/src/github.com/containers/podman-tui"
    config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/containers/podman-tui",
        type: "nfs",
        nfs_version: 4,
        nfs_udp: false

    config.vm.provider :libvirt do |domain|
        domain.memory = 4096
        domain.cpus = 2
    end

    install_go_env = <<-BASH
set -e
if [ ! -d "/usr/local/go" ]; then
    cd /tmp && wget https://golang.org/dl/go1.17.3.linux-amd64.tar.gz
    cd /usr/local
    tar xvzf /tmp/go1.17.3.linux-amd64.tar.gz
    echo 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.bashrc
    echo 'export GOROOT=/usr/local/go' >> /home/vagrant/.bashrc
    echo 'export GOBIN=/home/vagrant/go/bin' >> /home/vagrant/.bashrc
    echo 'export GOPRIVATE=github.com/containers/podman-tui' >> /home/vagrant/.bashrc
    echo 'export PATH=/usr/local/go/bin:$PATH:$GOPATH/bin' >> /home/vagrant/.bashrc
fi
export GOPATH=/home/vagrant/go
export GOBIN=/home/vagrant/go/bin
export GOROOT=/usr/local/go
export GOPRIVATE=github.com/containers/podman-tui
export PATH=/usr/local/go/bin:$PATH:$GOPATH/bin

BASH

    install_podman_env = <<-BASH
dnf -y install podman
dnf install -y btrfs-progs-devel device-mapper-devel gpgme-devel libassuan-devel
BASH

    config.vm.provision "shell", inline: "dnf -y install git-core wget gcc make bats tmux golint"
    config.vm.provision "shell", inline: install_go_env
    config.vm.provision "shell", inline: install_podman_env
    config.vm.provision "shell", inline: "chown -R vagrant.vagrant /home/vagrant/go"

end

