package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"reflect"
)

var Env Config

type GeneralConfig interface {
	Identify() string
	IsLocal() bool
}

type SpecificConfig interface {
	Identifier
	Typer
}

type Identifier interface {
	_identify() string
}

type Typer interface {
	_isLocal() bool
}

type Config struct {
	Domain    string `envconfig:"USPY_DOMAIN" required:"true" default:"localhost"`
	Port      string `envconfig:"USPY_PORT" required:"true" default:"8081"`
	JWTSecret string `envconfig:"USPY_JWT_SECRET" required:"true" default:"my_secret"`
	Mode      string `envconfig:"USPY_MODE" required:"true" default:"dev"`
	AESKey    string `envconfig:"USPY_AES_KEY" required:"true" default:"71deb5a48500599862d9e2170a60f90194a49fa81c24eacfe9da15cb76ba8b11"` // only used in dev
	RateLimit string `envconfig:"USPY_RATE_LIMIT" default:"5-S"`                                                                           // see github.com/ulule/limiter for more info

	Local  LocalConfig
	Remote RemoteConfig
}

func (c Config) IsLocal() bool {
	return !reflect.DeepEqual(c.Local, LocalConfig{})
}

func (c Config) Identify() string {
	if c.IsLocal() {
		return c.Local._identify()
	} else {
		return c.Remote._identify()
	}
}

func (c Config) Redact() Config {
	c.AESKey = "[REDACTED]"
	c.JWTSecret = "[REDACTED]"
	c.Domain = "[REDACTED]"
	c.Local.FirestoreKeyPath = "[REDACTED]"
	c.Remote.ProjectID = "[REDACTED]"
	return c
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("did not parse .env file, falling to default env variables")
	}

	if _, ok := os.LookupEnv("USPY_FIRESTORE_KEY"); ok {
		log.Println("Running backend locally")
		var lc LocalConfig
		envconfig.MustProcess("uspy", &lc)

		log.Printf("local env variables set: %#v\n", lc)
		Env.Local = lc
	} else if _, ok := os.LookupEnv("USPY_PROJECT_ID"); ok {
		log.Println("Running backend remotely")
		var rc RemoteConfig
		envconfig.MustProcess("uspy", &rc)

		log.Printf("remote env variables set: %#v\n", rc)
		Env.Remote = rc
	} else {
		log.Fatal("Could not initialize backend because neither the Firestore Key nor the Project ID were specified")
	}

	if err := envconfig.Process("uspy", &Env); err != nil {
		log.Fatal("could not process default env variables: ", err)
	}

	log.Printf("env variables set: %#v\n", Env.Redact())
}
