package vm

// dependency
// yum install libvirt-devel gcc
// brew install libvirt gcc
// apt install libvirt-dev gcc

import (
	"context"
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/libvirt/libvirt-go"
	"github.com/mohaijiang/computeshare-client/internal/conf"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//go:embed cloud-init.cfg.tmpl
var cloudInitTemp string

// VirtManager virtual machine management client
type VirtManager struct {
	conn    *libvirt.Connect
	log     *log.Helper
	cli     *client.Client
	workdir string

	noVncConnectionCancelMap map[string]func()
}

// NewVirtManager create virtManager
func NewVirtManager(logger log.Logger, cli *client.Client, data *conf.Data) (IVirtManager, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, err
	}
	var vmDir string

	fmt.Println("data: ", data.Workdir)
	if data.Workdir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		vmDir = path.Join(homeDir, "vm")
	} else {
		vmDir = data.Workdir
	}
	manager := &VirtManager{
		conn:                     conn,
		log:                      log.NewHelper(logger),
		workdir:                  vmDir,
		cli:                      cli,
		noVncConnectionCancelMap: make(map[string]func()),
	}
	return manager, err
}

func (v *VirtManager) initBaseData() {

	// md5sum f0432ad697f5762c28980a397c4e8d60
	// https://g.alpha.hamsternet.io/ipfs/QmZnCDgtSBQzHTyv2Ksku4zAxq9t7yUJwWGHUZAj2oX4AB?filename=ubuntu-20.04.qcow2.bak

	for _, item := range downloadFiles {
		err := v.DownloadFile(item)
		for {
			if err == nil {
				break
			}
			v.log.Error("下载镜像失败,重试")
			err = v.DownloadFile(item)
		}

	}

}

func (v *VirtManager) DownloadFile(image Image) error {

	stats, err := os.Open(path.Join(v.workdir, image.Filename))
	defer stats.Close()
	if err == nil {
		hash := md5.New()
		if _, err := io.Copy(hash, stats); err == nil {
			md5Hash := hash.Sum(nil)
			md5String := hex.EncodeToString(md5Hash)
			if md5String == image.MD5 {
				return nil
			}
		}
	}

	file, err := os.OpenFile(path.Join(v.workdir, image.Filename), os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	resp, err := http.Get(image.DownloadUrl)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (v *VirtManager) getCopyDiskFile(name string) string {
	return fmt.Sprintf("%s/%s.qcow2", v.workdir, name)
}

func (v *VirtManager) getBackupDiskFile(name string) string {
	return fmt.Sprintf("%s/%s.qcow2.backup", v.workdir, name)
}

func (v *VirtManager) getBaseImageName(image string) string {
	return downloadFiles[image].Filename
}

func (v *VirtManager) getBaseImagePath(image string) string {
	return path.Join(v.workdir, v.getBaseImageName(image))
}

// Create 创建虚拟机
func (v *VirtManager) Create(param *queueTaskV1.ComputeInstanceTaskParamVO) (string, error) {
	v.log.Info("start the virtual machine")

	imageInfo := downloadFiles[param.Image]

	if _, err := os.Stat(v.getCopyDiskFile(param.InstanceId)); errors.Is(err, os.ErrNotExist) {
		_ = os.MkdirAll(path.Dir(v.getCopyDiskFile(param.InstanceId)), os.ModePerm)

		fmt.Println("cp", v.getBaseImagePath(param.Image), v.getCopyDiskFile(param.InstanceId))
		cmd := exec.Command("cp", v.getBaseImagePath(param.Image), v.getCopyDiskFile(param.InstanceId))
		err := cmd.Run()
		if err != nil {
			fmt.Println("Execute Command failed:" + err.Error())
		}
	}

	err := v.generateCloudInitCfg(param.Name, param.GetPublicKey(), param.GetPassword(), param.DockerCompose)
	if err != nil {
		return "", err
	}

	// 实例化成cloud_init iso

	//cloud-localds cloud-init.iso cloud-init.cfg
	cmd := exec.Command("cloud-localds", "cloud-init.iso", "cloud-init.cfg")
	cmd.Dir = v.workdir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	fmt.Println(string(output))

	// virsh net-start default
	// virsh net-autostart default

	//virt-install \
	//--name $VM_NAME \
	//--memory 1024 \
	//--disk ubuntu-20.04.qcow2,device=disk,bus=virtio \
	//--disk cloud-init.iso,device=cdrom \
	//--os-type linux \
	//--os-variant ubuntu20.04 \
	//--virt-type kvm \
	//--graphics vnc,listen=0.0.0.0 \
	//--network network=default,model=virtio \
	//--noautoconsole \
	//--import

	vncPort := v.GetMaxVncPort() + 1
	cmds := []string{
		"virt-install",
		"--name", param.InstanceId,
		"--memory", strconv.Itoa(int(param.Memory * 1024)),
		"--vcpus", strconv.Itoa(int(param.Cpu)),
		"--disk", fmt.Sprintf("%s,device=disk,bus=virtio", v.getCopyDiskFile(param.InstanceId)),
		"--disk", "cloud-init.iso,device=cdrom",
		//"--os-type", imageInfo.OsType,
		"--os-variant", imageInfo.OsVariant,
		"--virt-type", "kvm",
		"--graphics", fmt.Sprintf("vnc,listen=0.0.0.0,port=%d", vncPort),
		"--network", "network=default,model=virtio",
		"--noautoconsole",
		"--import"}
	fmt.Println(cmds)
	cmd = exec.Command(cmds[0], cmds[1:]...)
	cmd.Dir = v.workdir
	output, err = cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("Execute Command failed:", err.Error())
		return "", err
	}

	ctx := context.Background()
	err = v.runNoVncCommandWithDocker(ctx, fmt.Sprintf("vnc_%s", param.InstanceId), v.GetVncWebsocketPort(param.InstanceId), int32(vncPort))

	return param.InstanceId, err
}

func (v *VirtManager) generateCloudInitCfg(name, publicKey, password string, dockercompose string) error {
	// 实例化初始cloud-init.iso
	tmpl, err := template.New("cloud-init").Funcs(template.FuncMap{"indent": indent}).Parse(cloudInitTemp)
	if err != nil {
		return err
	}

	// 创建或打开一个文件以便写入
	file, err := os.Create(path.Join(v.workdir, "cloud-init.cfg"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if password == "" {
		password = "123456"
	}

	data := CloudInitConf{
		Hostname:      name,
		Password:      password,
		PublicKey:     publicKey,
		DockerCompose: dockercompose,
	}

	// 使用模板将数据渲染并写入文件
	err = tmpl.Execute(file, data)
	return err
}

// Start start the virtual machine
func (v *VirtManager) Start(name string) error {

	d, err := v.conn.LookupDomainByName(name)
	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			v.log.Error("free libvirt.Domain fail")
		}
	}(d)

	if err != nil {
		return err
	}

	if active, err := d.IsActive(); active {
		return err
	}

	return d.Create()
}

// Stop shutdown the virtual machine
func (v *VirtManager) Stop(name string) error {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			v.log.Error("free libvirt.Domain fail")
		}
	}(d)

	return d.Shutdown()
}

// Reboot restart the virtual machine
func (v *VirtManager) Reboot(name string) error {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}

	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			v.log.Error("free libvirt.Domain fail")
		}
	}(d)

	return d.Reboot(libvirt.DOMAIN_REBOOT_DEFAULT)
}

// Shutdown the virtual machine
func (v *VirtManager) Shutdown(name string) error {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			v.log.Error("free libvirt.Domain fail")
		}
	}(d)

	return d.ShutdownFlags(libvirt.DOMAIN_SHUTDOWN_ACPI_POWER_BTN)
}

// Destroy the virtual machine
func (v *VirtManager) Destroy(instanceId string) error {
	d, err := v.conn.LookupDomainByName(instanceId)
	if err != nil {
		return err
	}
	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			v.log.Error("free libvirt.Domain fail")
		}
	}(d)

	state, _, err := d.GetState()
	if err != nil {
		return err
	}

	if libvirt.DOMAIN_SHUTOFF == state {
		return d.Undefine()
	}

	err = d.Destroy()
	if err != nil {
		return err
	}
	err = d.Undefine()

	if err != nil {
		return err
	}

	ctx := context.Background()
	_ = v.cli.ContainerRemove(ctx, fmt.Sprintf("vnc_%s", instanceId), types.ContainerRemoveOptions{Force: true})

	// 删除虚拟机文件
	return os.Rename(v.getCopyDiskFile(instanceId), v.getBackupDiskFile(instanceId))
}

// Status View status
func (v *VirtManager) Status(name string) (libvirt.DomainState, error) {
	dom, err := v.conn.LookupDomainByName(name)
	if err != nil {
		fmt.Println("err", err)
		return libvirt.DOMAIN_NOSTATE, nil
	}

	state, _, err := dom.GetState()

	//DOMAIN_NOSTATE     = DomainState(C.VIR_DOMAIN_NOSTATE)
	//	DOMAIN_RUNNING     = DomainState(C.VIR_DOMAIN_RUNNING)
	//	DOMAIN_BLOCKED     = DomainState(C.VIR_DOMAIN_BLOCKED)
	//	DOMAIN_PAUSED      = DomainState(C.VIR_DOMAIN_PAUSED)
	//	DOMAIN_SHUTDOWN    = DomainState(C.VIR_DOMAIN_SHUTDOWN)
	//	DOMAIN_CRASHED     = DomainState(C.VIR_DOMAIN_CRASHED)
	//	DOMAIN_PMSUSPENDED = DomainState(C.VIR_DOMAIN_PMSUSPENDED)
	//	DOMAIN_SHUTOFF     = DomainState(C.VIR_DOMAIN_SHUTOFF)

	defer func(dom *libvirt.Domain) {
		err := dom.Free()
		if err != nil {
			v.log.Error("free libvirt.Domain fail")
		}
	}(dom)
	return state, err
}

// GetIp get runtime ip
func (v *VirtManager) GetIp(name string) (string, error) {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return "", err
	}
	var dis []libvirt.DomainInterface
	failTimes := 0
	for {
		if failTimes > 180 {
			return "", err
		}
		dis, err = d.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE)
		if err != nil || len(dis) == 0 {
			failTimes++
			time.Sleep(time.Second)
			fmt.Println("fail time is :", failTimes)
			continue
		} else {
			fmt.Println("success time is :", failTimes)
			break
		}
	}

	for _, di := range dis {
		if len(di.Addrs) == 0 {
			continue
		}

		for _, ipAddress := range di.Addrs {
			return ipAddress.Addr, nil
		}
	}
	return "", errors.New("cannot get vm ip address")
}

// GetAccessPort get runtime port
func (v *VirtManager) GetAccessPort(name string) int {
	return 22
}

// GetVncPort get vnc port
func (v *VirtManager) GetVncPort(name string) int {
	d, err := v.conn.LookupDomainByName(name)
	if err != nil {
		return 0
	}

	domainXml, err := d.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	if err != nil {
		return 0
	}
	var libvirtDomainRoot LibvirtDomainRoot
	err = xml.Unmarshal([]byte(domainXml), &libvirtDomainRoot)
	if err != nil {
		return 0
	}
	fmt.Println(libvirtDomainRoot)
	if libvirtDomainRoot.Devices.Graphics.Type == "vnc" {
		return libvirtDomainRoot.Devices.Graphics.Port
	}

	return 0

}

func (v *VirtManager) GetMaxVncPort() int {
	defaultPort := 5900
	maxVncPort := defaultPort
	domains, err := v.conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		return defaultPort
	}
	fmt.Println(len(domains))

	for _, domain := range domains {
		name, err := domain.GetName()
		if err != nil {
			return defaultPort
		}
		port := v.GetVncPort(name)
		if port > maxVncPort {
			maxVncPort = port
		}
	}

	return maxVncPort
}

func (v *VirtManager) VncOpen(name string, publicVncPort int32) error {

	ctx, _ := context.WithCancel(context.Background())

	port := v.GetVncPort(name)
	return v.runNoVncCommandWithDocker(ctx, name, publicVncPort, int32(port))
}

func (v *VirtManager) VncClose(name string) error {
	ctx := context.Background()

	return v.stopNoVncCommandWithDocker(ctx, name)
}

func (v *VirtManager) runNoVncCommandWithDocker(ctx context.Context, containerName string, listenPort, vncPort int32) error {
	imageName := "hamstershare/novnc-websockify:latest"

	list, err := v.cli.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", imageName)),
	})
	if err != nil {
		return err
	}

	if len(list) == 0 {
		_, err = v.cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			return err
		}
	}

	containerConfig := &container.Config{
		Image: imageName,
		Cmd: []string{
			strconv.Itoa(int(listenPort)),
			fmt.Sprintf("%s:%d", GetLocalIP(), vncPort),
		},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: map[nat.Port]struct{}{
			nat.Port(fmt.Sprintf("%d/tcp", listenPort)): struct{}{},
		},
	}

	//hostConfig := &container.HostConfig{
	//	PortBindings: nat.PortMap{
	//		nat.Port(fmt.Sprintf("%d/tcp", listenPort)): []nat.PortBinding{
	//			{HostIP: "0.0.0.0", HostPort: strconv.Itoa(listenPort)},
	//		},
	//	},
	//}

	resp, err := v.cli.ContainerCreate(
		ctx,
		containerConfig,
		nil, //hostConfig,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		v.log.Error(err)
		return err
	}

	if err := v.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		v.log.Error(err)
		return err
	}
	return err
}

func (v *VirtManager) stopNoVncCommandWithDocker(ctx context.Context, containerName string) error {
	return v.cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{
		Force: true,
	})
}

func (v *VirtManager) runNoVncCommand(ctx context.Context, listenPort, vncPort int) {
	cmd := exec.CommandContext(ctx, "/snap/bin/novnc", "--listen", strconv.Itoa(listenPort), "--vnc", fmt.Sprintf("localhost:%d", vncPort))
	cmd.Dir = "/snap/novnc/current"
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		v.log.Info("Error starting command: ", err)
		return
	}

	pid := cmd.Process.Pid
	fmt.Println("pid:", pid)

	ppids := make([]int, 0)
	ccmd := exec.Command("/usr/bin/pgrep", "-P", fmt.Sprintf("%d", pid))

	output, err := ccmd.Output()
	if err != nil {
		fmt.Println("Error queue processes: ", err)
	} else {
		pidStrings := strings.Fields(string(output))

		for _, pidString := range pidStrings {
			fmt.Println("Child Process ID : ", pidString)
			ppid, _ := strconv.Atoi(pidString)
			ppids = append(ppids, ppid)
		}
	}

	err = cmd.Wait()

	if errors.Is(ctx.Err(), context.Canceled) {

		for _, ppid := range ppids {
			ccmd = exec.Command("kill", "-9", strconv.Itoa(ppid))
			ccmd.Output()
		}

		return
	}

	if err != nil {
		v.log.Info("Command failed: ", err)
		return
	}
}

func (v *VirtManager) GetVncWebsocketPort(name string) int32 {
	return 6800
}

func (v *VirtManager) GetVncWebsocketIP(instanceId string) (string, error) {
	ctx := context.Background()
	inspect, err := v.cli.ContainerInspect(ctx, fmt.Sprintf("vnc_%s", instanceId))
	if err != nil {
		return "", err
	}

	return inspect.NetworkSettings.IPAddress, nil
}

func indent(spaces int, s string) string {
	indentataion := strings.Repeat(" ", spaces)
	return indentataion + strings.ReplaceAll(s, "\n", "\n"+indentataion)
}
