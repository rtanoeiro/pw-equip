package equip

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Email string `json:"email"`
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".pw-equip-change")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// SaveEmail saves the email to the config file
func SaveEmail(email string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	config := Config{
		Email: email,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// LoadEmail loads the email from the config file
func LoadEmail() (string, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return "", err
	}

	// If config file doesn't exist, return empty string
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return "", err
	}

	return config.Email, nil
}
