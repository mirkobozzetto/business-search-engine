package tables

import (
	"csv-importer/api/middleware"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tableService TableService
}

func NewHandler(tableService TableService) *Handler {
	if tableService == nil {
		slog.Error("tableService is nil")
		os.Exit(1)
	}

	return &Handler{
		tableService: tableService,
	}
}

func (h *Handler) ListTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)

		tables, err := h.tableService.ListAllTables(c.Request.Context())
		if err != nil {
			slog.Error("failed to list tables",
				slog.String("error", err.Error()),
			)
			responseHelper.Error("failed to get tables: "+err.Error(), 500)
			return
		}

		responseHelper.Success(tables)
	}
}

func (h *Handler) GetTableInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)

		tableName := c.GetString("tableName")
		if tableName == "" {
			responseHelper.Error("table name not found", 400)
			return
		}

		tableInfo, err := h.tableService.GetTableInfo(c.Request.Context(), tableName)
		if err != nil {
			slog.Error("failed to get table info",
				slog.String("table", tableName),
				slog.String("error", err.Error()),
			)
			responseHelper.Error(err.Error(), 400)
			return
		}

		responseHelper.Success(tableInfo)
	}
}

func (h *Handler) GetTableColumns() gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)

		tableName := c.GetString("tableName")
		if tableName == "" {
			responseHelper.Error("table name not found", 400)
			return
		}

		columns, err := h.tableService.GetTableColumns(c.Request.Context(), tableName)
		if err != nil {
			slog.Error("failed to get table columns",
				slog.String("table", tableName),
				slog.String("error", err.Error()),
			)
			responseHelper.Error(err.Error(), 400)
			return
		}

		responseHelper.Success(columns)
	}
}

func (h *Handler) GetCompleteStructure() gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)

		structures, err := h.tableService.GetCompleteStructure(c.Request.Context())
		if err != nil {
			slog.Error("failed to get complete structure",
				slog.String("error", err.Error()),
			)
			responseHelper.Error("failed to get complete structure: "+err.Error(), 500)
			return
		}

		responseHelper.Success(structures)
	}
}
