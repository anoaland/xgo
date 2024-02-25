package database

import (
	"encoding/json"
	"log"
)

type SqlServerDatabaseConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

func (config SqlServerDatabaseConfig) JsonString() string {
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		log.Fatal("Can not marshal config: ", err)
	}

	return string(jsonConfig)
}

func (config SqlServerDatabaseConfig) IsValid() bool {
	notValid := config.Name == "" || config.Host == "" || config.User == "" || config.Password == "" || config.Port == ""

	return !notValid
}
