package database

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(config *PgDatabaseConfig, opts ...gorm.Option) *gorm.DB {

	dbname := config.Name
	host := config.Host
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host,
		config.User,
		config.Password,
		dbname,
		config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), opts...)

	if err != nil {
		db.Logger.Error(context.Background(), fmt.Sprintf("failed to connect database '%s' on '%s'", dbname, host))
		db.Logger.Error(context.Background(), err.Error())
		panic(err)
	}

	db.Logger.Info(context.Background(), fmt.Sprintf("Successfully connected to database '%s' on '%s'", dbname, host))

	return db
}
