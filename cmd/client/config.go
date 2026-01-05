package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Email - Email used to register the user
// BarChangeKey - Key used to change the skill bar
// TimingChange - Timing in milliseconds between each click
// ChangeSetKeyCode - Code of the key used to change the set
// ChangeSetKeyChar - Character of the key used to change the set
// Keys - Keys used to press the items in the set
type Config struct {
	Email               string   `json:"email"`
	Hwid                string   `json:"hwid"`
	BarChangeKey        string   `json:"barchangekey"`
	InBetweenTimeClicks int      `json:"inbetweentimeclicks"`
	ChangeSetKeyCode    uint16   `json:"changesetkey"`
	ChangeSetKeyChar    string   `json:"changesetkeychar"`
	Keys                []string `json:"keys"`
}

// getConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
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

func SaveConfig(email, hwid, keyChange, changeSetKeyChar string, inBetweenTimeClicks int, keys []string, changeSetKeyCode uint16) error {
	config := Config{
		Email:               email,
		Hwid:                hwid,
		BarChangeKey:        keyChange,
		InBetweenTimeClicks: inBetweenTimeClicks,
		Keys:                keys,
		ChangeSetKeyCode:    changeSetKeyCode,
		ChangeSetKeyChar:    changeSetKeyChar,
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("Successfully saved config into file")

	return os.WriteFile(configPath, data, 0644)
}

// LoadEmail loads the email from the config file
func LoadConfig() (Config, error) {
	configPath, err := GetConfigPath()
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
