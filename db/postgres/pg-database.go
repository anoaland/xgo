package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(config *PgDatabaseConfig) *gorm.DB {

	dbname := config.Name
	host := config.Host
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host,
		config.User,
		config.Password,
		dbname,
		config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database '%s' on '%s'", dbname, host)
	}

	log.Printf("Successfully connected to database '%s' on '%s'", dbname, host)

	return db
}
