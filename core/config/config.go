package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	FOLDER_NAME        = ".ado"
	FILENAME           = "config.ini"
	ENVIRONMENT_PREFIX = "ADO"
)

// Config represents the configuration structure
type Config struct {
	PAT string
	URL string
}

// GetConfigFilePath returns the path to the config file in a folder under the home directory
func GetConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, FOLDER_NAME)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("unable to create config directory: %w", err)
	}

	return filepath.Join(configDir, FILENAME), nil
}

// ReadConfig reads configuration from a file
func Read(profileName *string) (*Config, error) {
	configFile, err := GetConfigFilePath()
	if err != nil {
		log.Fatalf("Error getting config file path: %v", err)
	}

	profile := GetProfile(profileName)

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// // Automatically override values with environment variables if they exist
	// viper.AutomaticEnv()

	// // Set environment variable prefix to avoid conflicts
	// viper.SetEnvPrefix(ENVIRONMENT_PREFIX)

	// // Bind environment variables
	// viper.BindEnv("pat")
	// a := viper.GetString("PAT")
	// fmt.Println(">>", a)
	// // viper.BindEnv(profile + ".pat")
	// // viper.BindEnv(profile + ".org")

	var config Config
	config.URL = viper.GetString(profile + "url")
	config.PAT = viper.GetString(profile + "pat")

	return &config, nil
}

func GetProfile(profileName *string) string {
	profile := "default."
	if profileName != nil && *profileName != "" {
		profile = "profile " + *profileName + "."
	}
	return profile
}

// WriteConfig writes the configuration to a file
func Write(config *Config) error {

	configFile, err := GetConfigFilePath()
	if err != nil {
		log.Fatalf("Error getting config file path: %v", err)
	}

	viper.SetConfigType("ini")
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// p := "HELLO"
	// profile := GetProfile(&p)
	profile := GetProfile(nil)
	viper.Set(profile+"url", config.URL)
	viper.Set(profile+"pat", config.PAT)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	// configFile, err := GetConfigFilePath()
	// if err != nil {
	// 	log.Fatalf("Error getting config file path: %v", err)
	// }

	// // 	return WriteConfig(configFile, config)
	// // }

	// // // WriteConfig writes the configuration to a file
	// // func WriteConfig(configFile string, config *Config) error {
	// profile := GetProfile(nil)
	// viper.Set(profile+"pat", config.PAT)
	// viper.Set(profile+"org", config.Organization)

	// viper.SetConfigFile(configFile)
	// if err := viper.WriteConfig(); err != nil {
	// 	return fmt.Errorf("error writing config file: %w", err)
	// }

	return nil
}
