package config

// For a local configuration, you must have a path to your IAM key that can support firestore read/writes
type LocalConfig struct {
	FirestoreKeyPath string `envconfig:"USPY_FIRESTORE_KEY"`
}

func (lc LocalConfig) _identify() string {
	return lc.FirestoreKeyPath
}

func (LocalConfig) _isLocal() bool {
	return true
}
