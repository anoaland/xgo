package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository[M interface{}, D IDto[M, D], DList IDto[M, DList], DCreate ICreateDto[M], DUpdate IUpdateDto[M, D]] struct {
	db *gorm.DB
}

type BriefRepository[M interface{}, D IDto[M, D], DCreate ICreateDto[M]] struct {
	*Repository[M, D, D, DCreate, D]
}

func Brief[M interface{}, D IDto[M, D], DCreate ICreateDto[M]](context *gorm.DB) *BriefRepository[M, D, DCreate] {
	return &BriefRepository[M, D, DCreate]{
		Repository: &Repository[M, D, D, DCreate, D]{
			db: context,
		},
	}
}

func New[M interface{}, D IDto[M, D], DList IDto[M, DList], DCreate ICreateDto[M], DUpdate IUpdateDto[M, D]](context *gorm.DB) *Repository[M, D, DList, DCreate, DUpdate] {
	return &Repository[M, D, DList, DCreate, DUpdate]{
		db: context,
	}
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) Tx(tx *gorm.DB) *Repository[M, D, DList, DCreate, DUpdate] {
	return New[M, D, DList, DCreate, DUpdate](tx)
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) Create(payload DCreate) (*D, error) {
	values := payload.ToModel()
	err := r.db.Create(&values).Error
	if err != nil {
		return nil, err
	}

	var d D
	res := d.FromModel(values)
	return res, nil
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) CreateRaw(payload DCreate) (*M, error) {
	values := payload.ToModel()
	err := r.db.Create(&values).Error
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) Update(payload DUpdate, whereQuery interface{}, whereArgs ...interface{}) error {
	values := payload.ToModel()
	err := r.db.Model(&values).Where(whereQuery, whereArgs...).Updates(&values).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) SoftDelete(whereQuery interface{}, whereArgs ...interface{}) error {
	model := new(M)
	now := time.Now().UTC()
	err := r.db.Model(model).Where(whereQuery, whereArgs...).Update("DeletedAt", now).Error

	return err
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) FindAll(conds interface{}, orderBy interface{}, args ...interface{}) ([]DList, error) {
	model := new(*M)
	var rows *[]M
	r.db.Model(&model).Where(conds, args...).Order(orderBy).Find(&rows)

	return r.MapList(rows), nil
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) MapList(rows *[]M) []DList {
	results := make([]DList, 0, len(*rows))
	for _, row := range *rows {
		var d DList
		results = append(results, *d.FromModel(&row))
	}

	return results
}

func (r *Repository[M, D, DList, DCreate, DUpdate]) FindOne(conds ...interface{}) (*D, error) {
	var value *M
	result := r.db.First(&value, conds)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &NotFoundError{Message: "Record Not Found"}
	}

	var d D
	res := d.FromModel(value)
	return res, nil
}

type IDto[M interface{}, D interface{}] interface {
	ToModel() *M
	FromModel(*M) *D // IDto[M]
}

type IUpdateDto[M interface{}, D interface{}] interface {
	ToModel() *M
	// FromModel(M) D // IDto[M]
}

type ICreateDto[M interface{}] interface {
	ToModel() *M
}
