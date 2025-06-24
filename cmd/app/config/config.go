package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppFlags struct {
	LogFormat string
}

func ParseFlags() AppFlags {
	logFormat := flag.String("logf", "", "Logs format")
	flag.Parse()
	return AppFlags{
		LogFormat: *logFormat,
	}
}

func MustLoad(cfgPath string, cfg any) {
	if cfgPath == "" {
		log.Fatal("Config path is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist by this path: %s", cfgPath)
	}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}
}

type RabbitMQ struct {
	Hostname  string `yaml:"host"`
	Port      uint16 `yaml:"port"`
	QueueName string `yaml:"queue_name"`
}
type HTTPConfig struct {
	Address string `yaml:"address"`
}

type AppConfig struct {
	HTTPConfig `yaml:"rabbit_mq"`
	RabbitMQ   `yaml:"http"`
}
