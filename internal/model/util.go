package model

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Comparator struct {
	// represents way to compare threshold and result
	// 0: NA, 1: eq. 2: g, 3: ge, 4: l, 5: le
	Operator  int      `json:"operator"`
	Threshold string   `json:"threshold"`
	Arg       []string `json:"arg,omitempty"`
}

type Paginator struct {
	PageSize uint
	PageNum  uint
	Total    uint
	Rows     interface{}
	Filters  map[interface{}][]interface{}
	Pager    func(db *gorm.DB) *gorm.DB

	Err error
}

func (p *Paginator) ApplyQuery(db *gorm.DB, dest interface{}) *Paginator {
	for f := range p.Filters {
		db = db.Where(f, p.Filters[f]...)
	}

	var total int64
	err := db.Count(&total).Error
	p.Err = err
	if err == nil {
		p.Total = uint(total)
	}

	err = db.Scopes(p.Pager).Find(dest).Error
	p.Err = err
	if err == nil {
		p.Rows = dest
	}

	return p
}

func (p *Paginator) AddFilter(query interface{}, args ...interface{}) {
	if p.Filters == nil {
		p.Filters = make(map[interface{}][]interface{})
	}
	p.Filters[query] = args
}

// func (p *Paginator) ApplyFilters(db *gorm.DB) *gorm.DB {
// 	for f := range p.Filters {
// 		db = db.Where(f, p.Filters[f]...)
// 	}
// 	return db
// }

func (p *Paginator) SetPager(r *http.Request, order string) {
	p.Pager = func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Order(order).Offset(offset).Limit(pageSize)
	}
}
