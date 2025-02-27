package logger

import (
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func LogFromOpts(opts ...gorm.Option) gormlogger.Interface {
	var log gormlogger.Interface
	if len(opts) > 0 {
		cfg := opts[0].(*gorm.Config)
		if cfg != nil {
			log = cfg.Logger
		}
	}

	if log == nil {
		log = gormlogger.Default
	}

	return log.LogMode(gormlogger.Error)
}
