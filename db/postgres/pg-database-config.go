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

func (config PgDatabaseConfig) Dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
}
