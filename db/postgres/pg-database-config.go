package database

import (
	"encoding/json"
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

func IsValid(config PgDatabaseConfig) bool {
	notValid := config.Name == "" || config.Host == "" || config.User == "" || config.Password == "" || config.Port == ""

	return !notValid
}
