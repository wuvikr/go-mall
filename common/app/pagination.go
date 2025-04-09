package app

import (
	"go-mall/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

type pagination struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

func NewPagination(c *gin.Context) *pagination {
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize <= 0 {
		pageSize = config.App.Pagination.DefaultPageSize
	}
	if pageSize > config.App.Pagination.MaxPageSize {
		pageSize = config.App.Pagination.MaxPageSize
	}
	return &pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p *pagination) SetTotalRows(total int) {
	p.TotalRows = total
}

func (p *pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *pagination) GetPageSize() int {
	return p.PageSize
}

func (p *pagination) GetPage() int {
	return p.Page
}
