package config

type RemoteConfig struct {
	ProjectID string `envconfig:"USPY_PROJECT_ID"`
}

func (rc RemoteConfig) _identify() string {
	return rc.ProjectID
}

func (rc RemoteConfig) _isLocal() bool {
	return false
}
