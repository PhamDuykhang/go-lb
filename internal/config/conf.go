package config

const (
	DockerProvider = "docker"
	K8SProvider    = "k8s"

	DiscoveryLabel = "ktech.com.klb.group"
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
