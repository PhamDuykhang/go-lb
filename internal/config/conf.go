package config

const (
	DockerProvider = "docker"
	K8SProvider    = "k8s"

	DiscoveryLabel = "ktech.loadbalacing"
)

type (
	LBConfig struct {
		DiscoveryConfig Discovery
	}

	Discovery struct {
		Provider  string
		MathLabel string
	}
)
