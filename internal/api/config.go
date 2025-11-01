package api

import (
	"database/sql"
	"fmt"
	"log"
	"pw-equip-change/internal/database"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DefaultReadHeaderTimeout = 5 * time.Second
var DefaultMaxConns = 10
var DefaultConnMaxLifetime = time.Minute * 3

type Config struct {
	Environment   string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
	MySQLHost     string
	MySQLPort     string
	ApiPort       string
	DB            *database.Queries
}

type EquipConfig struct {
	Config
}

type MockConfig struct {
	Config
}

type ConfigLoader interface {
	LoadEquipConfig() *Config
}

func (config *EquipConfig) LoadEquipConfig() {
	*config = EquipConfig{
		Config: Config{
			Environment:   GetEnvVar("ENV", "dev"),
			MySQLUser:     GetEnvVar("MYSQL_USER", "pw-equip"),
			MySQLPassword: GetEnvVar("MYSQL_PASSWORD", "password"),
			MySQLDatabase: GetEnvVar("MYSQL_DATABASE", "equip"),
			MySQLHost:     GetEnvVar("MYSQL_HOST", "localhost"),
			MySQLPort:     GetEnvVar("MYSQL_PORT", "3307"),
			ApiPort:       GetEnvVar("API_PORT", "8989"),
		},
	}

	mySQLEquipDatabase := GetEnvVar("MYSQL_EQUIP_DATABASE", "equip")
	equipDBURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.MySQLUser,
		config.MySQLPassword,
		config.MySQLHost,
		config.MySQLPort,
		mySQLEquipDatabase,
	)
	log.Printf("Database URL: %s", equipDBURL)
	db, err := sql.Open("mysql", equipDBURL)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		panic(err)
	}

	db.SetConnMaxLifetime(DefaultConnMaxLifetime)
	db.SetMaxOpenConns(DefaultMaxConns)
	db.SetMaxIdleConns(DefaultMaxConns)
	config.DB = database.New(db)
}

// For mocking config, we don't load the database
func (config *MockConfig) LoadEquipConfig() *MockConfig {
	return &MockConfig{
		Config: Config{
			ApiPort:       GetEnvVar("API_PORT", "8989"),
			MySQLHost:     GetEnvVar("MYSQL_HOST", "localhost"),
			MySQLUser:     GetEnvVar("MYSQL_USER", "pw-server-tools"),
			MySQLPassword: GetEnvVar("MYSQL_PASSWORD", "password"),
			MySQLPort:     GetEnvVar("MYSQL_PORT", "3306"),
		},
	}
}
