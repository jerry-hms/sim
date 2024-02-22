package pagination

import "math"

func GetPagination() *Pagination {
	return &Pagination{
		Rows: new(interface{}),
	}
}

type Pagination struct {
	Page       int         `json:"page" form:"page"`
	PageSize   int         `json:"page_size" form:"page_size"`
	TotalRows  int64       `json:"total_rows" form:"total_rows"`
	TotalPages int         `json:"total_pages" form:"total_pages"`
	Rows       interface{} `json:"rows" form:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

func (p *Pagination) GetEnd() int {
	return p.Page * (p.PageSize - 1)
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) CountTotalPages(count int64) {
	p.TotalRows = count
	p.TotalPages = int(math.Ceil(float64(p.TotalRows) / float64(p.GetLimit())))
}
