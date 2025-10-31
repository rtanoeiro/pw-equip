package equip

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestLoadEquipConfig(t *testing.T) {
	config := Config{}
	config.LoadEquipConfig()
	if config.MySQLHost != "localhost" {
		t.Errorf("Expected MySQLHost to be localhost, got %s", config.MySQLHost)
	}
}
