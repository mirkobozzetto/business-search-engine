package middleware

import (
	"csv-importer/api/models"
	"time"

	"github.com/gin-gonic/gin"
)

type ResponseHelper struct {
	c *gin.Context
}

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("startTime", time.Now())
		c.Set("responseHelper", &ResponseHelper{c: c})
		c.Next()
	}
}

func GetResponseHelper(c *gin.Context) *ResponseHelper {
	helper, exists := c.Get("responseHelper")
	if !exists {
		return &ResponseHelper{c: c}
	}
	return helper.(*ResponseHelper)
}

func (r *ResponseHelper) Success(data any) {
	r.c.JSON(200, models.Success(data))
}

func (r *ResponseHelper) SuccessWithMeta(data any, meta models.Meta) {
	enrichedMeta := r.enrichMeta(meta)
	r.c.JSON(200, models.SuccessWithMeta(data, enrichedMeta))
}

func (r *ResponseHelper) Error(message string, statusCode int) {
	r.c.JSON(statusCode, models.Error(message))
}

func (r *ResponseHelper) enrichMeta(meta models.Meta) models.Meta {
	startTime, exists := r.c.Get("startTime")
	if exists {
		duration := time.Since(startTime.(time.Time))
		meta.Duration = duration.Milliseconds()
	}

	// Add pagination info
	if meta.Limit > 0 && meta.Count > 0 {
		meta.Page = (meta.Offset / meta.Limit) + 1
		if meta.Total > 0 {
			meta.Pages = (meta.Total + meta.Limit - 1) / meta.Limit
		}
	}

	return meta
}

func (r *ResponseHelper) SuccessWithPagination(data any, count, total, limit, offset int) {
	meta := models.Meta{
		Count:  count,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	r.SuccessWithMeta(data, meta)
}

// Helper to calculate total count when needed
func (r *ResponseHelper) GetTotalCount(tableName, columnName, searchValue string) int {
	// This would be implemented based on your needs
	// For now, returning 0 to avoid breaking existing functionality
	return 0
}
