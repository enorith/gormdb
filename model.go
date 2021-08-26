package gormdb

import "gorm.io/gorm"

type Model struct {
	*gorm.DB
}
