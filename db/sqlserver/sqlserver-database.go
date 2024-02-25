package database

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func Connect(config *SqlServerDatabaseConfig) *gorm.DB {
	dsn := config.Dsn(nil)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	dbname := config.Name
	host := config.Host
	if err != nil {
		log.Fatalf("failed to connect database '%s' on '%s'", dbname, host)
	}

	log.Printf("Successfully connected to database '%s' on '%s'", dbname, host)

	return db
}

func (config *SqlServerDatabaseConfig) Dsn(dbname *string) string {

	if dbname == nil {
		dbname = &config.Name
	}

	host := config.Host

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		config.User,
		config.Password,
		host,
		config.Port,
		*dbname,
	)

	return dsn
}
