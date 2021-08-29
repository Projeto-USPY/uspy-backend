package config

// For a remote configuration, you must provide the project id that contains the firestore database
type RemoteConfig struct {
	ProjectID string `envconfig:"USPY_PROJECT_ID"`
	Mailjet          // email verification is needed when running remotely
}

func (rc RemoteConfig) _identify() string {
	return rc.ProjectID
}

func (rc RemoteConfig) _isLocal() bool {
	return false
}
