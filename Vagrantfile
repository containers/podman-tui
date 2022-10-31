VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
    config.vm.hostname = "fedora36"
    config.vm.box = "fedora/36-cloud-base"
    config.vm.box_version = "36-20220504.1"

    config.vm.provision "shell", inline: "mkdir -p /home/vagrant/go"
    config.vm.synced_folder ".", "/home/vagrant/go/src/podman-tui",
        type: "nfs",
        nfs_version: 4,
        nfs_udp: false

    config.vm.provider :libvirt do |domain|
        domain.memory = 4096
        domain.cpus = 2
    end

    setup_env = <<-BASH
dnf -y update
dnf -y install glibc-static git-core wget gcc make bats tmux rpkg go-rpm-macros python3-pip
BASH

    setup_go = <<-BASH
dnf -y install golang golint
echo 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.bashrc
echo 'export GOBIN=/home/vagrant/go/bin' >> /home/vagrant/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> /home/vagrant/.bashrc
mkdir /home/vagrant/go/bin
BASH

    setup_podman = <<-BASH
dnf -y install podman
BASH

    config.vm.provision "shell", inline: setup_env
    config.vm.provision "shell", inline: setup_go
    config.vm.provision "shell", inline: setup_podman
    config.vm.provision "shell", inline: "chown -R vagrant.vagrant /home/vagrant/go"

end
