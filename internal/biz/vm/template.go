package vm

type CloudInitConf struct {
	InstanceId              string
	Password                string
	Hostname                string
	PublicKey               string
	DockerCompose           string
	PrometheusDockerCompose string
}

type SystemInfo struct {
	Hostname       string
	Arch           string
	TotalCpu       int32
	TotalMemory    int32
	OccupiedCpu    int32
	OccupiedMemory int32
}
