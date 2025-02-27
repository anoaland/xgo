package database

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func Connect(config *SqlServerDatabaseConfig, opts ...gorm.Option) *gorm.DB {
	dsn := config.Dsn(nil)
	db, err := gorm.Open(sqlserver.Open(dsn), opts...)

	dbname := config.Name
	host := config.Host
	if err != nil {
		db.Logger.Error(context.Background(), fmt.Sprintf("failed to connect database '%s' on '%s'", dbname, host))
		db.Logger.Error(context.Background(), err.Error())
		panic(err)
	}

	db.Logger.Info(context.Background(), fmt.Sprintf("Successfully connected to database '%s' on '%s'", dbname, host))

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
