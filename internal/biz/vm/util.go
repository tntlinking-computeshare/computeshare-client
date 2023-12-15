package vm

import (
	"encoding/xml"
	"net"
)

var downloadFiles map[string]Image

type Image struct {
	Name        string
	Filename    string
	DownloadUrl string
	MD5         string
	OsType      string
	OsVariant   string
}

func init() {
	downloadFiles = make(map[string]Image)
	downloadFiles["ubuntu:20.04"] = Image{
		Name:        "ubuntu-20.04",
		Filename:    "ubuntu-20.04.qcow2.bak",
		DownloadUrl: "https://g.alpha.hamsternet.io/ipfs/QmZnCDgtSBQzHTyv2Ksku4zAxq9t7yUJwWGHUZAj2oX4AB?filename=ubuntu-20.04.qcow2.bak",
		MD5:         "f0432ad697f5762c28980a397c4e8d60",
		OsType:      "linux",
		OsVariant:   "ubuntu20.04",
	}
}

type LibvirtDomainRoot struct {
	XMLNAME xml.Name             `xml:"domain"`
	Devices LibVirtDomainDevices `xml:"devices"`
}

type LibVirtDomainDevices struct {
	XMLName  xml.Name              `xml:"devices"`
	Graphics LibVirtDomainGraphics `xml:"graphics"`
}

type LibVirtDomainGraphics struct {
	XMLName xml.Name `xml:"graphics"`
	Type    string   `xml:"type,attr"`
	Port    int      `xml:"port,attr"`
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
