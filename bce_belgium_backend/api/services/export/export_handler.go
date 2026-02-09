package export

import (
	"csv-importer/api/models"
	"encoding/csv"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	exportService ExportService
}

func NewHandler(exportService ExportService) *Handler {
	if exportService == nil {
		slog.Error("exportService is nil")
		os.Exit(1)
	}

	return &Handler{
		exportService: exportService,
	}
}

func (h *Handler) ExportData() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		columnName := c.Query("column")
		searchValue := c.Query("search")
		limitStr := c.DefaultQuery("limit", "10000")
		format := c.DefaultQuery("format", "csv")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100000 {
			c.JSON(400, models.Error("invalid limit parameter (max 100000)"))
			return
		}

		opts := ExportOptions{
			TableName:   tableName,
			ColumnName:  columnName,
			SearchValue: searchValue,
			Limit:       limit,
			Format:      format,
		}

		result, err := h.exportService.PrepareExportData(c.Request.Context(), opts)
		if err != nil {
			slog.Error("failed to prepare export data",
				slog.String("table", tableName),
				slog.String("column", columnName),
				slog.String("search", searchValue),
				slog.Int("limit", limit),
				slog.String("format", format),
				slog.String("error", err.Error()),
			)
			c.JSON(500, models.Error("export failed: "+err.Error()))
			return
		}

		if format == "json" {
			h.handleJSONExport(c, opts, result)
		} else {
			h.handleCSVExport(c, opts, result)
		}
	}
}

func (h *Handler) handleJSONExport(c *gin.Context, opts ExportOptions, result *ExportResult) {
	response := map[string]any{
		"table": opts.TableName,
		"data":  result.Data,
		"meta":  result.Meta,
	}

	c.JSON(200, models.Success(response))
}

func (h *Handler) handleCSVExport(c *gin.Context, opts ExportOptions, result *ExportResult) {
	filename := h.exportService.GenerateFilename(opts) + ".csv"

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	if err := writer.Write(result.Columns); err != nil {
		slog.Error("failed to write CSV header",
			slog.String("filename", filename),
			slog.String("error", err.Error()),
		)
		c.JSON(500, models.Error("failed to write CSV header"))
		return
	}

	if len(result.Data) == 0 {
		writer.Write([]string{"No data found"})
		c.Status(http.StatusOK)
		return
	}

	for _, row := range result.Data {
		record := make([]string, len(result.Columns))
		for i, col := range result.Columns {
			if val, exists := row[col]; exists && val != nil {
				record[i] = fmt.Sprintf("%v", val)
			} else {
				record[i] = ""
			}
		}

		if err := writer.Write(record); err != nil {
			slog.Error("failed to write CSV row",
				slog.String("filename", filename),
				slog.String("error", err.Error()),
			)
			break
		}
	}

	c.Status(http.StatusOK)
}
