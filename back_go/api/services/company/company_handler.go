package company

import (
	"csv-importer/api/models"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	companyService CompanyService
}

func NewHandler(companyService CompanyService) *Handler {
	if companyService == nil {
		slog.Error("companyService is nil")
		os.Exit(1)
	}

	return &Handler{
		companyService: companyService,
	}
}

func (h *Handler) SearchByNaceCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		naceCode := c.Query("code")
		limitStr := c.DefaultQuery("limit", "50")

		if naceCode == "" {
			c.JSON(400, models.Error("nace code parameter 'code' is required"))
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		result, err := h.companyService.SearchByNaceCode(c.Request.Context(), naceCode, limit)
		if err != nil {
			slog.Error("failed to search by nace code",
				slog.String("nace_code", naceCode),
				slog.Int("limit", limit),
				slog.String("error", err.Error()),
			)
			c.JSON(500, models.Error("search failed: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}

func (h *Handler) SearchByDenomination() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		limitStr := c.DefaultQuery("limit", "50")

		if query == "" {
			c.JSON(400, models.Error("denomination query parameter 'q' is required"))
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		result, err := h.companyService.SearchByDenomination(c.Request.Context(), query, limit)
		if err != nil {
			slog.Error("failed to search by denomination",
				"query", query,
				"limit", limit,
				"error", err.Error(),
			)
			c.JSON(500, models.Error("search failed: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}

func (h *Handler) SearchByZipcode() gin.HandlerFunc {
	return func(c *gin.Context) {
		zipcode := c.Query("q")
		limitStr := c.DefaultQuery("limit", "50")

		if zipcode == "" {
			c.JSON(400, models.Error("zipcode query parameter 'q' is required"))
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		result, err := h.companyService.SearchByZipcode(c.Request.Context(), zipcode, limit)
		if err != nil {
			slog.Error("failed to search by zipcode",
				"zipcode", zipcode,
				"limit", limit,
				"error", err.Error(),
			)
			c.JSON(500, models.Error("search failed: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}

func (h *Handler) SearchMultiCriteria() gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "50")

		criteria := models.CompanySearchCriteria{
			NaceCode:     c.Query("nace"),
			Denomination: c.Query("denomination"),
			ZipCode:      c.Query("zipcode"),
			Status:       c.Query("status"),
		}

		if criteria.NaceCode == "" && criteria.Denomination == "" && criteria.ZipCode == "" && criteria.Status == "" {
			c.JSON(400, models.Error("at least one search criteria required (nace, denomination, zipcode, status)"))
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 1000 {
			c.JSON(400, models.Error("invalid limit parameter"))
			return
		}

		result, err := h.companyService.SearchMultiCriteria(c.Request.Context(), criteria, limit)
		if err != nil {
			slog.Error("failed to search multi criteria",
				"criteria", fmt.Sprintf("%+v", criteria),
				"limit", limit,
				"error", err.Error(),
			)
			c.JSON(500, models.Error("search failed: "+err.Error()))
			return
		}

		c.JSON(200, models.Success(result))
	}
}
