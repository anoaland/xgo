package database

import (
	"context"
	"fmt"
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

	if strings.HasPrefix(sql, "SELECT") {
		r.Interface.Trace(ctx, begin, fc, err)
		return
	}

	if code == 0 {
		if !strings.HasPrefix(sql, "SELECT") && !strings.Contains(sql, "pg_catalog.pg_description") {
			// r.Interface.Trace(ctx, begin, fc, err)
			log.Println("==============================", sql)
			r.Statements = append(r.Statements, sql)
		}
	}
}

// PrintAutoMigrateSql prints the SQL statements for auto migrating the given model structs.
//
// db: the database connection
// dst: the model structs to auto migrate
// string: the SQL statements for auto migrating the given model structs
func PrintAutoMigrateSqlx(db *gorm.DB, dst ...interface{}) string {

	db.Config.PrepareStmt = true

	// thanks to: https://stackoverflow.com/a/66246127/1586914
	recorder := RecorderLogger{logger.Default.LogMode(logger.Silent), []string{}}
	session := db.Session(
		&gorm.Session{
			DryRun: true,
			Logger: &recorder,
		})
	// err := session.AutoMigrate(dst...)

	migrator := session.Dialector.Migrator(session)
	err := migrator.AutoMigrate(dst...)

	if err != nil {
		log.Fatalf("failed to generate automigrate sql: %v", err)
	}

	sql := strings.Join(recorder.Statements, ";\r\n")

	return sql + ";"
}

// thanks to: https://github.com/go-gorm/gorm/issues/3851#issuecomment-929752108
func PrintAutoMigrateSql(db *gorm.DB, dst ...interface{}) string {
	tx := db.Begin()
	var statements []string
	tx.Callback().Raw().Register("record_migration", func(tx *gorm.DB) {
		statements = append(statements, tx.Statement.SQL.String())
	})
	if err := tx.AutoMigrate(dst...); err != nil {
		panic(err)
	}
	tx.Rollback()
	tx.Callback().Raw().Remove("record_migration")
	for _, s := range statements {
		fmt.Println(s)
	}

	return strings.Join(statements, ";\r\n") + ";"
}
