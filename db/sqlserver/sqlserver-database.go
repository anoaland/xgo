package database

import (
	"context"
	"fmt"

	"github.com/anoaland/xgo/db/logger"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func Connect(config *SqlServerDatabaseConfig, opts ...gorm.Option) *gorm.DB {
	log := logger.LogFromOpts(opts...)

	dsn := config.Dsn(nil)
	dbname := config.Name
	host := config.Host

	db, err := gorm.Open(sqlserver.Open(dsn), opts...)
	if err != nil {
		log.Error(context.Background(), fmt.Sprintf("failed to connect database '%s' on '%s'", dbname, host))
		log.Error(context.Background(), err.Error())
		panic(err)
	}

	log.Info(context.Background(), fmt.Sprintf("Successfully connected to database '%s' on '%s'", dbname, host))

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
