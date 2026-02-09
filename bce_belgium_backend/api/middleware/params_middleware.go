package middleware

import (
	"csv-importer/api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ParsedParams struct {
	Limit  int
	Offset int
	Format string
}

func ParseLimitParam(defaultLimit, maxLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", strconv.Itoa(defaultLimit))
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > maxLimit {
			c.JSON(400, models.Error("invalid limit parameter (max "+strconv.Itoa(maxLimit)+")"))
			c.Abort()
			return
		}

		c.Set("limit", limit)
		c.Next()
	}
}

func ParseOffsetParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		offsetStr := c.DefaultQuery("offset", "0")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(400, models.Error("invalid offset parameter"))
			c.Abort()
			return
		}

		c.Set("offset", offset)
		c.Next()
	}
}

func ParseFormatParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		format := c.DefaultQuery("format", "json")
		if format != "json" && format != "csv" {
			c.JSON(400, models.Error("invalid format parameter (json or csv)"))
			c.Abort()
			return
		}

		c.Set("format", format)
		c.Next()
	}
}

func ParseSortParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		sortBy := c.DefaultQuery("sort", "")
		order := c.DefaultQuery("order", "asc")

		if order != "asc" && order != "desc" {
			c.JSON(400, models.Error("invalid order parameter (asc or desc)"))
			c.Abort()
			return
		}

		c.Set("sortBy", sortBy)
		c.Set("order", order)
		c.Next()
	}
}

func GetParsedParams(c *gin.Context) ParsedParams {
	return ParsedParams{
		Limit:  c.GetInt("limit"),
		Offset: c.GetInt("offset"),
		Format: c.GetString("format"),
	}
}
