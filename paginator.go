package gormdb

import (
	"math"

	"gorm.io/gorm"
)

var ResultFormater = func(data interface{}, meta PageMeta) map[string]interface{} {
	return map[string]interface{}{
		"data": data,
		"meta": meta,
	}
}

type Paginator struct {
	page, perPage int
	manger        *Manager
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

	conn, e := p.manger.GetConnection()
	if e != nil {
		return nil, e
	}
	var meta PageMeta
	meta.Page = p.page
	meta.PerPage = p.perPage
	meta.From = p.perPage*(p.page-1) + 1

	tx.Statement.Dest = targets
	session := tx.Session(&gorm.Session{})

	e = conn.Table("(?) as `aggregate`", session).Count(&meta.Total).Error

	if e != nil {
		return nil, e
	}
	meta.LastPage = int(math.Ceil(float64(meta.Total) / float64(p.perPage)))

	db := tx.Offset(p.perPage * (p.page - 1)).Limit(p.perPage).Find(targets)

	meta.To = meta.From + int(db.RowsAffected-1)

	return ResultFormater(targets, meta), nil
}

func NewPaginator(page, perPage int, ms ...*Manager) *Paginator {
	if page < 1 {
		page = 1
	}
	var manager *Manager
	if len(ms) > 0 {
		manager = ms[0]
	} else {
		manager = DefaultManager
	}

	return &Paginator{page: page, perPage: perPage, manger: manager}
}
