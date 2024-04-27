package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository[M interface{}, D IDto[M, D], DList IDto[M, DList], DCreate ICreateDto[M], DUpdate IUpdateDto[M, D]] struct {
	db *gorm.DB
}

func New[M interface{}, D IDto[M, D], DList IDto[M, DList], DCreate ICreateDto[M], DUpdate IUpdateDto[M, D]](context *gorm.DB) *Repository[M, D, DList, DCreate, DUpdate] {
	return &Repository[M, D, DList, DCreate, DUpdate]{
		db: context,
	}
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) Tx(tx *gorm.DB) *Repository[M, D, DList, DCreate, DUpdate] {
	return New[M, D, DList, DCreate, DUpdate](tx)
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) Create(payload DCreate) (D, error) {
	values := payload.ToModel()
	r.db.Create(&values)
	var d D
	return d.FromModel(values), nil
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) Update(payload DUpdate, whereQuery interface{}, whereArgs ...interface{}) (*D, error) {
	values := payload.ToModel()
	err := r.db.Model(&values).Where(whereQuery, whereArgs...).Updates(&values).Error
	if err != nil {
		return nil, err
	}

	res := payload.FromModel(values)
	return &res, nil
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) SoftDelete(model M, whereQuery interface{}, whereArgs ...interface{}) error {
	now := time.Now().UTC()
	err := r.db.Model(&model).Where(whereQuery, whereArgs...).Updates(
		map[string]interface{}{
			"deleted_at": &now,
		},
	).Error
	return err
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) FindAll(conds ...interface{}) ([]DList, error) {
	var rows *[]M
	r.db.Find(&rows, conds)

	return r.MapList(rows), nil
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) MapList(rows *[]M) []DList {
	results := make([]DList, 0, len(*rows))
	for _, row := range *rows {
		var d DList
		results = append(results, d.FromModel(row))
	}

	return results
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) FindOne(conds ...interface{}) (D, error) {
	var value *M
	result := r.db.First(&value, conds)

	var d D
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return d, &NotFoundError{Message: "Record Not Found"}
	}

	return d.FromModel(*value), nil
}

type IDto[M interface{}, D interface{}] interface {
	ToModel() M
	FromModel(M) D // IDto[M]
}

type IUpdateDto[M interface{}, D interface{}] interface {
	ToModel() M
	FromModel(M) D // IDto[M]
}

type ICreateDto[M interface{}] interface {
	ToModel() M
}
