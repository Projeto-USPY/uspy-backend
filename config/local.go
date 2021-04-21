package config

type LocalConfig struct {
	FirestoreKeyPath string `envconfig:"USPY_FIRESTORE_KEY"`
}

func (lc LocalConfig) _identify() string {
	return lc.FirestoreKeyPath
}

func (LocalConfig) _isLocal() bool {
	return true
}
