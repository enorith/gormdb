package gormdb

import "gorm.io/gorm"

type Builder[T interface{}] struct {
	*gorm.DB
}

func (b *Builder[T]) Query(fn func(*gorm.DB) *gorm.DB) *Builder[T] {
	return NewBuilder[T](b.DB.Scopes(fn))
}

func (b *Builder[T]) Get() (result []T, err error) {
	err = b.DB.Find(&result).Error
	return
}

func (b *Builder[T]) First(conds ...interface{}) (result T, err error) {
	err = b.DB.First(&result, conds...).Error
	return
}

func NewBuilder[T interface{}](tx *gorm.DB) *Builder[T] {
	return &Builder[T]{DB: tx}
}
