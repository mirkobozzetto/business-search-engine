package data

import (
	"csv-importer/api/models"
	"log/slog"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	dataService DataService
}

// NewHandler creates a new data handler with dependency injection
func NewHandler(dataService DataService) *Handler {
	if dataService == nil {
		slog.Error("dataService is nil")
		os.Exit(1)
	}

	return &Handler{
		dataService: dataService,
	}
}

func (h *Handler) PreviewTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		limitStr := c.DefaultQuery("limit", "5")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		result, err := h.dataService.PreviewTable(c.Request.Context(), tableName, limit)
		if err != nil {
			slog.Error("failed to preview table",
				slog.String("table", tableName),
				slog.Int("limit", limit),
				slog.String("error", err.Error()),
			)
			c.JSON(500, models.Error("failed to preview table: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}

func (h *Handler) GetColumnValues() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		limitStr := c.DefaultQuery("limit", "20")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		result, err := h.dataService.GetColumnValues(c.Request.Context(), tableName, columnName, limit)
		if err != nil {
			slog.Error("failed to get column values",
				slog.String("table", tableName),
				slog.String("column", columnName),
				slog.Int("limit", limit),
				slog.String("error", err.Error()),
			)
			c.JSON(500, models.Error("failed to get column values: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}
