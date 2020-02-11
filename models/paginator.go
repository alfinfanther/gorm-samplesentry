package models

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Pagination struct {
	Limit        int
	Page         int
	Last_ID      int
	Order        string
	TotalRecords int
	TotalPages   int64
	Db           *gorm.DB
}

func (p *Pagination) Paginate(c *gin.Context, db *gorm.DB) (*gorm.DB, error) {
	p.Db = db
	limitQuery := c.DefaultQuery("limit", "25")
	pageQuery := c.DefaultQuery("page", "1")
	lastIdQuery := c.Query("last_id")
	p.Order = c.DefaultQuery("order", "desc")

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		return db, errors.New("invalid parameter.")
	}
	p.Limit = int(math.Max(1, math.Min(10000, float64(limit))))

	if lastIdQuery != "" {
		// pagination 1
		last_id, err := strconv.Atoi(lastIdQuery)
		if err != nil {
			return db, errors.New("invalid parameter.")
		}
		p.Last_ID = int(math.Max(0, float64(last_id)))
		if p.Order == "asc" {
			return db.Where("id > ?", p.Last_ID).Limit(p.Limit).Order("id asc"), nil
		} else {
			return db.Where("id < ?", p.Last_ID).Limit(p.Limit).Order("id desc"), nil
		}
	}

	// pagination 2
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		return db, errors.New("invalid parameter.")
	}
	p.Page = int(math.Max(1, float64(page)))
	return db.Offset(limit * (p.Page - 1)).Limit(p.Limit), nil
}

func (p *Pagination) CountRecords(tableName string, mapKv map[string]string, args []interface{}) {
	var count int
	//db.Count(&count)
	if mapKv["join_conditions"] != "" {
		p.Db.Table(tableName).Where(mapKv["query_conditions"], args...).Joins(mapKv["join_conditions"]).Count(&count)
	} else {
		p.Db.Table(tableName).Where(mapKv["query_conditions"], args...).Count(&count)
	}
	p.TotalRecords = count
	p.TotalPages = p.getTotalPages(strconv.Itoa(p.Limit), count)
}

func (p *Pagination) SetHeaderLink(c *gin.Context, index int) {
	var link string
	if p.Last_ID != 0 {
		link = fmt.Sprintf("<http://%v%v?limit=%v&last_id=%v&order=%v>; rel=\"next\"", c.Request.Host, c.Request.URL.Path, p.Limit, index, p.Order)
	} else {
		if p.Page == 1 {
			link = fmt.Sprintf("<http://%v%v?limit=%v&page=%v>; rel=\"next\"", c.Request.Host, c.Request.URL.Path, p.Limit, p.Page+1)
		} else {
			link = fmt.Sprintf("<http://%v%v?limit=%v&page=%v>; rel=\"next\",<http://%v%v?limit=%v&page=%v>; rel=\"prev\"", c.Request.Host, c.Request.URL.Path, p.Limit, p.Page+1, c.Request.Host, c.Request.URL.Path, p.Limit, p.Page-1)
		}
	}
	c.Header("Link", link)
	c.Header("NumRow", strconv.Itoa(p.TotalRecords))
	c.Header("NumPage", strconv.Itoa(int(p.TotalPages)))
}

func (p *Pagination) getTotalPages(perPage string, totalRecords int) int64 {
	perPageInt, _ := strconv.ParseInt(perPage, 10, 32)
	totalPages := float64(totalRecords) / float64(perPageInt)
	return int64(float64(totalPages) + float64(1.0))
}
