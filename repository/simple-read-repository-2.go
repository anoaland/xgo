package repository

import (
	"errors"

	"github.com/anoaland/xgo"
	"gorm.io/gorm"
)

type SimpleReadRepository2[M interface{}, D SimpleReadIDto2[M, D]] struct {
	db *gorm.DB
}

func NewSimpleReadRepo2[M interface{}, D SimpleReadIDto2[M, D]](context *gorm.DB) *SimpleReadRepository2[M, D] {
	return &SimpleReadRepository2[M, D]{
		db: context,
	}
}

func (r *SimpleReadRepository2[M, D]) findAll(conds *string, orderBy *string, args ...interface{}) ([]D, error) {
	model := new(*M)
	var rows *[]M

	q := r.db.Model(&model)

	if conds != nil {
		q = q.Where(conds, args...)
	}

	if orderBy != nil {
		q = q.Order(orderBy)
	}

	err := q.Find(&rows).Error

	if err != nil {
		return nil, xgo.NewHttpInternalError("E_SIMPLE_READ_REPO_FIND_ALL", err)
	}

	return r.MapList(rows), nil
}

func (r *SimpleReadRepository2[M, D]) FindAll() ([]D, error) {
	return r.findAll(nil, nil)
}

func (r *SimpleReadRepository2[M, D]) FindAllWithOrder(orderBy string) ([]D, error) {
	return r.findAll(nil, &orderBy)
}

func (r *SimpleReadRepository2[M, D]) FindAllWithConditionAndOrder(conds string, orderBy string, args ...interface{}) ([]D, error) {
	return r.findAll(&conds, &orderBy, args...)
}

func (r *SimpleReadRepository2[M, D]) FindAllWithCondition(conds string, args ...interface{}) ([]D, error) {
	return r.findAll(&conds, nil, args...)
}

func (r *SimpleReadRepository2[M, D]) MapList(rows *[]M) []D {
	results := make([]D, 0, len(*rows))
	for _, row := range *rows {
		var d D
		results = append(results, *d.FromModel(&row))
	}

	return results
}

func (r *SimpleReadRepository2[M, D]) FindOne(conds ...interface{}) (*D, error) {
	var value *M
	result := r.db.First(&value, conds)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &NotFoundError{Message: "Record Not Found"}
	}

	var d D
	res := d.FromModel(value)
	return res, nil
}

type SimpleReadIDto2[M interface{}, D interface{}] interface {
	FromModel(*M) *D // IDto[M]
}
