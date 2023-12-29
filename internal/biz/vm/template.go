package vm

type CloudInitConf struct {
	InstanceId              string
	Password                string
	Hostname                string
	PublicKey               string
	DockerCompose           string
	PrometheusDockerCompose string
}
