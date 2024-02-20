package service

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type JSONResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type PagedDataResp struct {
	JSONResp
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func getPageQuery(c *gin.Context) (int, int, error) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return 0, 0, err
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		return 0, 0, err
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return page, pageSize, nil
}
