package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// It represents the configuration
type Config struct {
	RabbitMQ RabbitMQConfig
	SMTP     SMTPConfig
}

// It represent RabbitMQ configuration
type RabbitMQConfig struct {
	URL string `mapstructure:"url"`
}

// It represent the SMTP provider configuration
type SMTPConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	URL      string `mapstructure:"url"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// It loads configuration from files and environment variables
func LoadConfig() (*Config, error) {
	viper.SetConfigName(fmt.Sprintf("config.%s", getEnv()))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/") // config.yaml file path
	viper.AutomaticEnv()             // Also reads from ENV

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading the configuration file: %w", err)
	}

	// Parse config file
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error parsing the configuration: %w", err)
	}

	log.Println("Configuration loaded successfully")
	return &config, nil
}

func getEnv() string {
	env_str := "APP_ENV"

	viper.BindEnv(env_str)
	env := viper.GetString(env_str) // Read the APP_ENV environment variable

	if env == "" {
		log.Printf("APP_ENV empty. Using default.")
		env = "dev" // Default to dev if not set
	}

	log.Printf("Using ENV [%s]", env)

	return env
}
