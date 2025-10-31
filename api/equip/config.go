package equip

import (
	"database/sql"
	"fmt"
	"log"
	"pw-equip-change/database"
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

type MockConfig struct {
	Config
}

type ConfigLoader interface {
	LoadEquipConfig() *Config
}

func (config *Config) LoadEquipConfig() *Config {
	config = &Config{
		Environment:   GetEnvVar("ENV", "dev"),
		MySQLUser:     GetEnvVar("MYSQL_USER", "pw-server-tools"),
		MySQLPassword: GetEnvVar("MYSQL_PASSWORD", "password"),
		MySQLDatabase: GetEnvVar("MYSQL_DATABASE", "pw-server-tools"),
		MySQLHost:     GetEnvVar("MYSQL_HOST", "localhost"),
		MySQLPort:     GetEnvVar("MYSQL_PORT", "3306"),
		ApiPort:       GetEnvVar("API_PORT", "8989"),
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
	return config
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
