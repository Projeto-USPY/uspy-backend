package config

// For a remote configuration, you must provide the project id that contains the firestore database
type RemoteConfig struct {
	ProjectID string `envconfig:"USPY_PROJECT_ID"`
}

func (rc RemoteConfig) _identify() string {
	return rc.ProjectID
}

func (rc RemoteConfig) _isLocal() bool {
	return false
}
