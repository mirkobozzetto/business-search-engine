package middleware

import (
	"csv-importer/api/helpers"
	"csv-importer/api/models"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func ValidateTableName() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("table")
		if tableName == "" {
			tableName = c.Param("name")
		}

		if err := helpers.ValidateTableName(tableName); err != nil {
			c.JSON(400, models.Error(err.Error()))
			c.Abort()
			return
		}

		c.Set("tableName", tableName)
		c.Next()
	}
}

func ValidateColumnName(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.GetString("tableName")
		if tableName == "" {
			c.JSON(500, models.Error("table name not found in context"))
			c.Abort()
			return
		}

		columnName := c.Param("column")
		if columnName == "" {
			c.Next()
			return
		}

		if err := helpers.ValidateColumnExists(db, tableName, columnName); err != nil {
			c.JSON(400, models.Error(err.Error()))
			c.Abort()
			return
		}

		c.Set("columnName", columnName)
		c.Next()
	}
}

func ValidateSearchQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		searchValue := c.Query("q")
		if searchValue == "" {
			c.JSON(400, models.Error("search query 'q' is required"))
			c.Abort()
			return
		}

		c.Set("searchValue", searchValue)
		c.Next()
	}
}
