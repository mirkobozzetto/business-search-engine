package search

import (
	"csv-importer/api/models"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	searchService SearchService
}

func NewHandler(searchService SearchService) *Handler {
	if searchService == nil {
		slog.Error("searchService is nil")
		os.Exit(1)
	}

	return &Handler{
		searchService: searchService,
	}
}

func (h *Handler) SearchTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		searchValue := c.Query("q")
		limitStr := c.DefaultQuery("limit", "50")

		if searchValue == "" {
			c.JSON(400, models.Error("search query 'q' is required"))
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		result, err := h.searchService.SearchInColumn(c.Request.Context(), tableName, columnName, searchValue, limit)
		if err != nil {
			slog.Error("failed to search in column",
				slog.String("table", tableName),
				slog.String("column", columnName),
				slog.String("query", searchValue),
				slog.Int("limit", limit),
				slog.String("error", err.Error()),
			)
			c.JSON(500, models.Error("search failed: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}

func (h *Handler) CountRows() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Param("column")
		searchValue := c.Query("q")

		if searchValue == "" {
			c.JSON(400, models.Error("search query 'q' is required"))
			return
		}

		result, err := h.searchService.CountMatches(c.Request.Context(), tableName, columnName, searchValue)
		if err != nil {
			slog.Error("failed to count matches",
				slog.String("table", tableName),
				slog.String("column", columnName),
				slog.String("query", searchValue),
				slog.String("error", err.Error()),
			)
			c.JSON(500, models.Error("count failed: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}

func (h *Handler) SearchNaceCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		searchValue := c.Query("q")
		limitStr := c.Query("limit")

		limit := parseOptionalLimit(limitStr, 0)

		result, err := h.searchService.SearchNaceCode(c.Request.Context(), searchValue, limit)
		if err != nil {
			slog.Error("failed to search NACE codes",
				slog.String("query", searchValue),
				slog.Int("limit", limit),
				slog.String("error", err.Error()),
			)
			c.JSON(500, models.Error("NACE search failed: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}

func (h *Handler) SearchMultipleColumns() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		searchValue := c.Query("q")
		columnsStr := c.Query("columns") // "col1,col2,col3"
		limitStr := c.DefaultQuery("limit", "50")

		if searchValue == "" {
			c.JSON(400, models.Error("search query 'q' is required"))
			return
		}

		if columnsStr == "" {
			c.JSON(400, models.Error("columns parameter is required"))
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		columns := parseColumns(columnsStr)
		if len(columns) == 0 {
			c.JSON(400, models.Error("at least one column is required"))
			return
		}

		result, err := h.searchService.SearchMultipleColumns(c.Request.Context(), tableName, columns, searchValue, limit)
		if err != nil {
			slog.Error("failed to search multiple columns",
				slog.String("table", tableName),
				slog.Any("columns", columns),
				slog.String("query", searchValue),
				slog.Int("limit", limit),
				slog.String("error", err.Error()),
			)
			c.JSON(500, models.Error("multi-column search failed: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}

func parseColumns(columnsStr string) []string {
	if columnsStr == "" {
		return nil
	}

	var columns []string
	for _, col := range strings.Split(columnsStr, ",") {
		col = strings.TrimSpace(col)
		if col != "" {
			columns = append(columns, col)
		}
	}

	return columns
}
