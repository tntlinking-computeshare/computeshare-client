package vm

var downloadFiles map[string]Image

type Image struct {
	Name        string
	Filename    string
	DownloadUrl string
	MD5         string
}

func init() {
	downloadFiles = make(map[string]Image)
	downloadFiles["ubuntu-20.04"] = Image{
		Name:        "ubuntu-20.04",
		Filename:    "ubuntu-20.04.qcow2.bak",
		DownloadUrl: "https://g.alpha.hamsternet.io/ipfs/QmZnCDgtSBQzHTyv2Ksku4zAxq9t7yUJwWGHUZAj2oX4AB?filename=ubuntu-20.04.qcow2.bak",
		MD5:         "f0432ad697f5762c28980a397c4e8d60",
	}
}
