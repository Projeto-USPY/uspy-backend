package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Env is the default variable for the environment to be loaded
var Env Config

// GeneralConfig is the interface for an environment configuration. It has only one method that must identify the type of authentication used
type GeneralConfig interface {
	Identify() string
}

// Config is the default configuration object, for more info see README.md
type Config struct {
	Domain    string `envconfig:"USPY_DOMAIN" required:"true" default:"localhost"`
	Port      string `envconfig:"USPY_PORT" required:"true" default:"8080"` // careful with this because cloud run must run on port 8080
	JWTSecret string `envconfig:"USPY_JWT_SECRET" required:"true" default:"my_secret"`
	Mode      string `envconfig:"USPY_MODE" required:"true" default:"local"`
	AESKey    string `envconfig:"USPY_AES_KEY" required:"true" default:"71deb5a48500599862d9e2170a60f90194a49fa81c24eacfe9da15cb76ba8b11"` // only used in dev
	RateLimit string `envconfig:"USPY_RATE_LIMIT"`                                                                                         // see github.com/ulule/limiter for more info

	FirestoreKeyPath  string `envconfig:"USPY_FIRESTORE_KEY"`
	MockFirestoreData bool   `envconfig:"USPY_MOCK_FIRESTORE_DATA" default:"false"`

	ProjectID string `envconfig:"USPY_PROJECT_ID"`

	Mailjet // email verification is needed in production
}

// IsUsingKey returns whether a firestore key is being used to authenticate with Firestore
func (c Config) IsUsingKey() bool {
	return c.FirestoreKeyPath != ""
}

// IsUsingProjectID returns whether the GCP Project ID is being used to authenticate with Firestore
func (c Config) IsUsingProjectID() bool {
	return c.ProjectID != ""
}

// Identify returns the type of the authentication being used by the configuration object
func (c Config) Identify() string {
	if c.IsUsingKey() {
		return c.FirestoreKeyPath
	}

	return c.ProjectID
}

// IsDev returns whether the configuration environment is in development mode
func (c Config) IsDev() bool {
	return c.Mode == "dev"
}

// IsLocal returns whether the configuration environment is in local mode
func (c Config) IsLocal() bool {
	return c.Mode == "local"
}

// IsProd returns whether the configuration environment is in production mode
func (c Config) IsProd() bool {
	return c.Mode == "prod"
}

// Redact can be used to print the environment config without exposing secret
func (c Config) Redact() Config {
	c.AESKey = "[REDACTED]"
	c.JWTSecret = "[REDACTED]"
	c.Domain = "[REDACTED]"
	c.FirestoreKeyPath = "[REDACTED]"
	c.ProjectID = "[REDACTED]"
	c.Mailjet.APIKey = "[REDACTED]"
	c.Mailjet.Secret = "[REDACTED]"
	return c
}

// TestSetup is used by the emulator, it will only load required defaults, no project-related identifiers
func TestSetup() {
	if (Env != Config{}) { // idempotent function
		return
	}

	if err := envconfig.Process("uspy", &Env); err != nil {
		log.Fatal("could not process default env variables: ", err)
	}

	log.Infof("env variables set %#v", Env)
}

// Setup parses the .env file (or uses defaults) to determine environment constants and variables
func Setup() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Info("did not parse .env file, falling to default env variables")
	}

	if err := envconfig.Process("uspy", &Env); err != nil {
		log.Fatal("could not process default env variables: ", err)
	}

	log.Info("env variables set", Env.Redact())

	if Env.IsUsingKey() {
		log.Info("Running backend with firestore key")

		if !utils.CheckFileExists(Env.FirestoreKeyPath) {
			log.Fatal("Could not find firestore key path: ", Env.FirestoreKeyPath)
		}
	} else if Env.IsUsingProjectID() {
		log.Info("Running backend with project ID")

		// setup email client
		Env.Mailjet.Setup()
	} else {
		log.Fatal("Could not initialize backend because neither the Firestore Key nor the Project ID were specified")
	}

}
