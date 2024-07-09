package repository

import (
	"math"
	"strings"
	"time"

	"github.com/anoaland/xgo/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SetKeywordLikeVarsByTotalExpr(keyword string, total int) (vars []interface{}) {
	for i := 0; i < total; i++ {
		vars = append(vars, "%"+strings.ToLower(keyword)+"%")
	}

	return vars
}

func LowerLikeQuery(field string) string {
	return "LOWER(" + field + ") LIKE ?"
}

func FilterPaginate(DB *gorm.DB, modelName interface{}, pagination *dto.Pagination, clauses []clause.Expression, joins []string) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	tx := DB.Model(modelName).Clauses(clauses...)

	if len(joins) > 0 {
		tx = txWithJoins(tx, joins)
	}

	tx.Count(&totalRows)

	pagination.TotalData = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {

		if len(joins) > 0 {
			db = txWithJoins(db, joins)
		}

		return db.Clauses(clauses...).Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func txWithJoins(tx *gorm.DB, joinList []string) *gorm.DB {
	for _, join := range joinList {
		tx.Joins(join)
	}

	return tx
}

func txWithPeloads(tx *gorm.DB, preloadList []string) *gorm.DB {
	for _, join := range preloadList {
		tx.Preload(join)
	}

	return tx
}

// special function if column didn't save with UTC time and you have to used this function to get the local time
func LocalTime(value time.Time) time.Time {

	_, offset := time.Now().Zone()

	offsetHours := (offset / 3600.0) * -1

	return value.UTC().Add(time.Duration(offsetHours) * time.Hour)
}
