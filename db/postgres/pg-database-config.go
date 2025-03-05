package database

import (
	"encoding/json"
	"fmt"
	"log"
)

type PgDatabaseConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

func (config PgDatabaseConfig) JsonString() string {
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		log.Fatal("Can not marshal config: ", err)
	}

	return string(jsonConfig)
}

func (config PgDatabaseConfig) IsValid() bool {
	notValid := config.Name == "" || config.Host == "" || config.User == "" || config.Password == "" || config.Port == ""

	return !notValid
}

// Dsn returns a fully qualified PostgreSQL connection string including database name.
// This is used for normal database operations.
func (config PgDatabaseConfig) Dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
}

// Dsn returns a fully qualified PostgreSQL connection string excluding the database name.
// This is used for testing operations like dropping or creating a test database.
func (config PgDatabaseConfig) DsnWithoutDB() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)
}
