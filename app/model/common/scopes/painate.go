package scopes

import (
	"gorm.io/gorm"
	"math"
	"sim/app/util/pagination"
)

func Paginate(p *pagination.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var totalRows int64
		db.Debug().Count(&totalRows)
		p.TotalRows = totalRows
		totalPages := int(math.Ceil(float64(totalRows) / float64(p.GetLimit())))
		p.TotalPages = totalPages

		return db.Offset(p.GetOffset()).Limit(p.GetLimit())
	}
}
