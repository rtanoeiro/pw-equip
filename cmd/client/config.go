package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Email string `json:"email"`
	BarChangeKey string `json:"barchangekey"`
	TimingChange string `json:"timingchange"`
	ChangeSetKeyCode uint16 `json:"changesetkey"`
	ChangeSetKeyChar string `json:"changesetkeychar"`
	Keys []string `json:"keys"`
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

func SaveConfig(email, keyChange, timingChange,changeSetKeyChar string, keys []string, changeSetKeyCode uint16) error {
	config := Config{
		Email: email,
		BarChangeKey: keyChange,
		TimingChange: timingChange,
		Keys: keys,
		ChangeSetKeyCode: changeSetKeyCode,
		ChangeSetKeyChar: changeSetKeyChar,
		
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}


	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("Successfullt saved config into file")

	return os.WriteFile(configPath, data, 0644)
}

// LoadEmail loads the email from the config file
func LoadConfig() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	// If config file doesn't exist, return empty string
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Config{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
