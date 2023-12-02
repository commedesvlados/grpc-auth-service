package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"
	EnvironmentLocal       = "local"
)

type Envs struct {
	Database struct {
		Path string `env:"PATH" env-required:"true"`
		//Host     string `env:"HOST" env-required:"true"`
		//Port     int    `env:"PORT" env-required:"true"`
		//User     string `env:"USER" env-required:"true"`
		//Password string `env:"PASSWORD" env-required:"true"`
		//Database string `env:"DB" env-required:"true"`
	} `env-prefix:"DATABASE_"`
	GRPC struct {
		Port    int           `env:"PORT" env-required:"true"`
		Timeout time.Duration `env:"TIMEOUT" env-required:"true"`
	} `env-prefix:"GRPC_"`
	TokenTTL time.Duration `env:"TOKEN_TTL" env-required:"true"`
}

var E *Envs
var onceE sync.Once

func ReadEnv(envPath string) {
	onceE.Do(func() {
		if err := godotenv.Load(envPath); err != nil {
			log.Fatalf("can't loading env variables, err: %s\n", err.Error())
		}

		log.Printf("[Config] Read environment variables, path: %s\n", envPath)
		E = &Envs{}
		if err := cleanenv.ReadEnv(E); err != nil {
			help, _ := cleanenv.GetDescription(E, nil)
			log.Println(help)
			log.Fatalln(err)
		}
	})
}

type Config struct {
	Env string `yaml:"environment" env-required:"true"` // local / development / production
}

var C *Config
var onceC sync.Once

func ReadConfig(configPath string) {
	onceC.Do(func() {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file does not exists: %s, err: %s\n", configPath, err.Error())
		}

		log.Printf("[Config] Read configuration variables, path: %s\n", configPath)
		C = &Config{}
		if err := cleanenv.ReadConfig(configPath, C); err != nil {
			help, _ := cleanenv.GetDescription(E, nil)
			log.Println(help)
			log.Fatalln(err)
		}
	})
}

func MustLoadVariables() {
	var fenv string
	flag.StringVar(&fenv, "env", "production/ development / local", "project environment")
	flag.Parse()

	envPath := ".env." + fenv
	configPath := fmt.Sprintf("config/config.%s.yaml", fenv)

	MustLoadVariablesByPath(envPath, configPath)
}

func MustLoadVariablesByPath(envPath, configPath string) {
	ReadEnv(envPath)
	ReadConfig(configPath)
}
