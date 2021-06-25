package plugins

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Cast struct {
}

func (c Cast) Name() string {
	return "cast"
}

func (c Cast) Initialize(db *gorm.DB) error {
	query := db.Callback().Query()

	query.After("gorm:after_query").Register("cast", func(db *gorm.DB) {
		for _, field := range db.Statement.Schema.Fields {
			str := field.Tag.Get("gorm")
			tags := strings.Split(str, ";")
			for _, tag := range tags {
				parts := strings.Split(tag, ":")
				if len(parts) > 1 && parts[0] == "cast" {
					field.DataType = "string"
				}
			}
			fmt.Println(field.TagSettings, field.DataType)
		}
	})

	return nil
}

type JsonField struct {
	field interface{}
}

func (j *JsonField) Value() (driver.Value, error) {
	return nil, nil
}

func (j JsonField) Scan(src interface{}) error {
	return nil
}
