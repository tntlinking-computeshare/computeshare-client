# 共享算力客户端安装-ubuntu

## virt

```shell
sudo su

## 安装libvirt
apt update
apt install qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils virtinst virt-manager libvirt-dev cloud-image-utils libvirt-dev gcc -y

vim /etc/libvirt/qemu.conf
## 使用 root 启动虚拟机
## user = "root"
systemctl restart libvirtd


## 安装docker

curl -fsSL https://get.docker.com -o install-docker.sh
sh install-docker.sh --mirror Aliyun

```
