package database

import (
	"context"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RecorderLogger struct {
	logger.Interface
	Statements []string
}

func (r *RecorderLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, code := fc()

	if code == 0 {
		if !strings.Contains(sql, "pg_catalog.pg_description") {
			r.Statements = append(r.Statements, sql)
		}
	}
}

// PrintAutoMigrateSql prints the SQL statements for auto migrating the given model structs.
//
// db: the database connection
// dst: the model structs to auto migrate
// string: the SQL statements for auto migrating the given model structs
func PrintAutoMigrateSql(db *gorm.DB, dst ...interface{}) string {
	// thanks to: https://stackoverflow.com/a/66246127/1586914

	recorder := RecorderLogger{logger.Default.LogMode(logger.Silent), []string{}}
	session := db.Session(&gorm.Session{DryRun: true, Logger: &recorder})
	err := session.AutoMigrate(dst...)
	if err != nil {
		log.Fatalf("failed to generate automigrate sql: %v", err)
	}

	sql := strings.Join(recorder.Statements, ";\r\n")

	return sql
}
