package naf

import (
	"net/http"
	"sirene-importer/api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *nafService
}

func NewHandler(service *nafService) *Handler {
	return &Handler{service: service}
}

func parseLimit(c *gin.Context, defaultLimit int) int {
	limit := defaultLimit
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > 10000 {
		limit = 10000
	}
	return limit
}

func parseOffset(c *gin.Context) int {
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			return parsed
		}
	}
	return 0
}

func (h *Handler) SearchByLabel(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, models.Error("q parameter required"))
		return
	}
	limit := parseLimit(c, 100)
	offset := parseOffset(c)
	codes, total, err := h.service.SearchByLabel(c.Request.Context(), q, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}
	pages := 0
	if limit > 0 && total > 0 {
		pages = (total + limit - 1) / limit
	}
	c.JSON(http.StatusOK, models.SuccessWithMeta(codes, models.Meta{
		Total:  total,
		Count:  len(codes),
		Limit:  limit,
		Offset: offset,
		Page:   page,
		Pages:  pages,
	}))
}

func (h *Handler) ListSections(c *gin.Context) {
	sections, err := h.service.ListSections(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(sections))
}

func (h *Handler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, models.Error("code parameter required"))
		return
	}
	nafCode, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	if nafCode == nil {
		c.JSON(http.StatusNotFound, models.Error("NAF code not found"))
		return
	}
	c.JSON(http.StatusOK, models.Success(nafCode))
}

func (h *Handler) GetBySection(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, models.Error("code parameter required"))
		return
	}
	codes, err := h.service.GetBySection(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(codes))
}
