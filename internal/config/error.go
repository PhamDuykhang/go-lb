package config

type (
	ForwardError struct {
		SourceURL string `json:"source_url"`
		Status    string `json:"status"`
		Message   string `json:"message"`
	}
)
