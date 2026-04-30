package helper

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mahimtalukder/oshudh-somadhan/internal/response"
)

func ParsePageLimit(c *gin.Context) (int, int, bool) {
	page := 1
	limit := 10

	if pageQuery := c.Query("page"); pageQuery != "" {
		value, err := strconv.Atoi(pageQuery)
		if value < 1 || err != nil {
			response.Error(c, http.StatusBadGateway, "Invalid page number", "Page number must be positive integer")
			return page, limit, false
		}
		page = value
	}

	if limitQuery := c.Query("limit"); limitQuery != "" {
		value, err := strconv.Atoi(limitQuery)
		if value < 1 || value > 100 || err != nil {
			response.Error(c, http.StatusBadGateway, "Invalid limit", "limit must be between 1 and 100")
			return page, limit, false
		}
		limit = value
	}

	return page, limit, true
}

func ParseOptionalIntQuery(c *gin.Context, key string) (*int, bool) {

	if valueString := c.Query(key); valueString != "" {
		value, err := strconv.Atoi(valueString)
		if err != nil || value < 1 {
			response.Error(c, http.StatusBadGateway, "Invalid query parameter", key+" must be a positive number")
			return nil, false
		}
		return &value, true
	}

	return nil, true
}
