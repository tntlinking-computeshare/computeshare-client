# 共享算力客户端安装-ubuntu

## libvirt
* 安装libvirt， 并将libvirt 启动用户修改为root

```shell
sudo su
apt update
apt install qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils virtinst virt-manager libvirt-dev cloud-image-utils libvirt-dev gcc -y

vim /etc/libvirt/qemu.conf
## 使用 root 启动虚拟机
## user = "root"
systemctl restart libvirtd
```

* 安装docker
```shell
curl -fsSL https://get.docker.com -o install-docker.sh
sh install-docker.sh --mirror Aliyun
```

## 安装配置computeshare-client客户端程序

```shell
##  创建工作目录
mkdir -p /var/lib/computeshare
mkdir -p /root/vm

cd /var/lib/computeshare

curl -L https://github.com/tntlinking-computeshare/computeshare-client/releases/download/0.1.0/computeshare-client_linux_amd64 -o computeshare-client 

cat > /var/lib/computeshare/config.yaml << EOF
server:
  http:
    addr: 0.0.0.0:18000
    timeout: 10s
  grpc:
    addr: 0.0.0.0:19000
    timeout: 10s
  p2p:
    gateway_ip: 61.172.179.73
    gateway_port: 7000
data:
  database:
    driver: mysql
    source: root:123456@tcp(127.0.0.1:3306)/test
  redis:
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s
  ipfs:
    url: 127.0.0.1:5001
  computer_power_api: "https://api.computeshare.newtouch.com"
  workdir: /root/vm
EOF

## 创建环境变量文件
cat /var/lib/computeshare/computeshare-client.env << EOF
HOME=/root
EOF

## 创建系统启动文件
cat > /lib/systemd/system/computeshare-client.service << EOF
[Unit]

Description=computeshare-client

[Service]
WorkingDirectory=/var/lib/computeshare
PIDFile=/run/computeshare-client.pid

EnvironmentFile=/var/lib/computeshare/computeshare-client.env

ExecStart=/var/lib/computeshare/computeshare-client

ExecReload=/bin/kill -SIGHUP $MAINPID

ExecStop=/bin/kill -SIGINT $MAINPID

[Install]

WantedBy=multi-user.target
EOF


## 启动算力提供者节点
systemctl daemon-reload && systemctl start computeshare-client
```
