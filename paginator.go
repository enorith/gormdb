package gormdb

import (
	"math"

	"gorm.io/gorm"
)

var ResultFormater = func(data interface{}, meta PageMeta) map[string]interface{} {
	return map[string]interface{}{
		"meta": meta,
		"data": data,
	}
}

type Paginator struct {
	page, perPage int
}

type PageMeta struct {
	Total    int64 `json:"total"`
	PerPage  int   `json:"per_page"`
	Page     int   `json:"page"`
	LastPage int   `json:"last_page"`
	From     int   `json:"from"`
	To       int   `json:"to"`
}

func (p *Paginator) Paginate(tx *gorm.DB, targets interface{}) (map[string]interface{}, error) {
	var meta PageMeta
	meta.Page = p.page
	meta.PerPage = p.perPage
	meta.From = p.perPage*(p.page-1) + 1
	newTx := tx.Session(&gorm.Session{
		NewDB: true,
	})
	tx.Statement.Dest = targets

	e := newTx.Table("(?) as `aggragate`", tx).Count(&meta.Total).Error
	if e != nil {
		return nil, e
	}

	db := tx.Limit(int(p.perPage)).Offset(int(p.perPage * (p.page - 1))).Find(targets)
	e = db.Error
	if e != nil {
		return nil, e
	}

	meta.LastPage = int(math.Ceil(float64(meta.Total) / float64(p.perPage)))

	meta.To = meta.From + int(db.RowsAffected-1)

	return ResultFormater(targets, meta), nil
}

func NewPaginator(page, perPage int) *Paginator {
	if page < 1 {
		page = 1
	}

	return &Paginator{page: page, perPage: perPage}
}
