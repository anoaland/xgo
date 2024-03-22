package repository

import (
	"errors"

	"github.com/anoaland/xgo"
	"gorm.io/gorm"
)

type SimpleReadRepository[M interface{}, D SimpleReadIDto[M, D]] struct {
	db *gorm.DB
}

func NewSimpleReadRepo[M interface{}, D SimpleReadIDto[M, D]](context *gorm.DB) *SimpleReadRepository[M, D] {
	return &SimpleReadRepository[M, D]{
		db: context,
	}
}

func (r *SimpleReadRepository[M, D]) FindAll(conds ...interface{}) ([]D, error) {
	var rows *[]M
	if err := r.db.Find(&rows, conds).Error; err != nil {
		return nil, xgo.NewHttpInternalError("E_SIMPLE_READ_REPO_FIND_ALL", err)
	}

	return r.MapList(rows), nil
}

func (r *SimpleReadRepository[M, D]) MapList(rows *[]M) []D {
	results := make([]D, 0, len(*rows))
	for _, row := range *rows {
		var d D
		results = append(results, d.FromModel(row))
	}

	return results
}

func (r *SimpleReadRepository[M, D]) FindOne(conds ...interface{}) (D, error) {
	var value *M
	result := r.db.First(&value, conds)

	var d D
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return d, &NotFoundError{Message: "Record Not Found"}
	}

	return d.FromModel(*value), nil
}

type SimpleReadIDto[M interface{}, D interface{}] interface {
	FromModel(M) D // IDto[M]
}
