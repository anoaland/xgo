package database

import (
	"context"
	"fmt"

	"github.com/anoaland/xgo/db/logger"
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

	log := logger.LogFromOpts(opts...)
	db, err := gorm.Open(postgres.Open(dsn), opts...)

	if err != nil {
		log.Error(context.Background(), "failed to connect database '%s' on '%s'", dbname, host)
		log.Error(context.Background(), err.Error())
		panic(err)
	}

	log.Info(context.Background(), "Successfully connected to database '%s' on '%s'", dbname, host)

	return db
}
