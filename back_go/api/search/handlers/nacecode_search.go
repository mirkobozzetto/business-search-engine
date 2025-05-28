package handlers

import (
	helperutils "csv-importer/api/helpers/utils"
	"csv-importer/api/models"
	searchutils "csv-importer/api/search/utils"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SearchNaceCode(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		searchValue := c.Query("q")
		limitStr := c.Query("limit")

		limit := searchutils.ParseOptionalLimit(limitStr, 0)

		query, args := searchutils.BuildNaceCodeQuery(searchValue, limit)

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(500, models.Error("database error: "+err.Error()))
			return
		}
		defer rows.Close()

		columns := []string{"nacecode", "activités", "libellé_fr", "omschrijving_nl"}
		data, err := helperutils.ScanRowsToMaps(rows, columns)
		if err != nil {
			c.JSON(500, models.Error("scan error: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(map[string]any{
			"query":   searchValue,
			"results": data,
			"meta": models.Meta{
				Count: len(data),
				Limit: limit,
				Total: len(data),
			},
		}))
	}
}
