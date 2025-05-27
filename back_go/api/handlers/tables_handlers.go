package handlers

import (
	"csv-importer/api/middleware"
	"csv-importer/api/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func ListTables(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)
		tableService := services.NewTableService(db)

		tables, err := tableService.ListAllTables()
		if err != nil {
			responseHelper.Error("failed to get tables: "+err.Error(), 500)
			return
		}

		responseHelper.Success(tables)
	}
}

func GetTableInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)
		tableService := services.NewTableService(db)

		tableName := c.GetString("tableName") // Set by middleware
		if tableName == "" {
			responseHelper.Error("table name not found", 400)
			return
		}

		tableInfo, err := tableService.GetTableInfo(tableName)
		if err != nil {
			responseHelper.Error(err.Error(), 400)
			return
		}

		responseHelper.Success(tableInfo)
	}
}

func GetTableColumns(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)
		tableService := services.NewTableService(db)

		tableName := c.GetString("tableName") // Set by middleware

		columns, err := tableService.GetTableColumns(tableName)
		if err != nil {
			responseHelper.Error(err.Error(), 400)
			return
		}

		responseHelper.Success(columns)
	}
}

func GetCompleteStructure(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)
		tableService := services.NewTableService(db)

		structures, err := tableService.GetCompleteStructure()
		if err != nil {
			responseHelper.Error("failed to get complete structure: "+err.Error(), 500)
			return
		}

		responseHelper.Success(structures)
	}
}
