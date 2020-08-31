package dicovery

const (
	Change ServiceChangeType = "update"
	Die    ServiceChangeType = "die"
	RePawn ServiceChangeType = "" // later
	Start  ServiceChangeType = "start"
)

type ServiceChangeType string

type (
	ServiceMetadata struct {
		ServiceID   string
		ServiceName string
		Action      ServiceChangeType
		NameAddress string
		Port        string
	}
)
