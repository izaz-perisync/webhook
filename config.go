package webhook

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   `yaml:"server"`
	JWT      `yaml:"jwt"`
	Database `yaml:"dataBase"`
	S3Key    `yaml:"s3Key"`
	Discord  `yaml:"discord"`
}

type Server struct {
	Port           int      `yaml:"port"`
	ExposedHeaders []string `yaml:"exposedHeaders"`
	AllowedHeaders []string `yaml:"allowedHeaders"`
	AllowedMethods []string `yaml:"allowedMethods"`
	AllowedOrigins []string `yaml:"allowedOrigins"`
}

type JWT struct {
	Key string `yaml:"key"`
}

type Database struct {
	DB string `yaml:"db"`
}

type S3Key struct {
	Region          string `yaml:"region"`
	Bucket          string `yaml:"bucket"`
	AccessKey       string `yaml:"accessKey"`
	SecretAccessKey string `yaml:"secretAccessKey"`
}

type Discord struct {
	WebhookUrl string `yaml:"webhookUrl"`
}

func GlobalConfig() *Config {

	config := Config{}

	path := "./"
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		// If there is an error reading the configuration file, panic with the error.
		log.Panic(err)
	}
	// Unmarshal the configuration file into the Config struct.
	if err := viper.Unmarshal(&config); err != nil {
		log.Panic(err)
	}
	// Return a pointer to the loaded Config struct.
	return &config
}
