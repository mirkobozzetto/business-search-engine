package company

import (
	"net/http"
	"strconv"
	"sirene-importer/api/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *companyService
}

func NewHandler(service *companyService) *Handler {
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

func (h *Handler) SearchByNafCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, models.Error("code parameter required"))
		return
	}
	limit := parseLimit(c, 100)
	result, err := h.service.SearchByNafCode(c.Request.Context(), code, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}

func (h *Handler) SearchByDenomination(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.Error("q parameter required"))
		return
	}
	limit := parseLimit(c, 100)
	result, err := h.service.SearchByDenomination(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}

func (h *Handler) SearchByCodePostal(c *gin.Context) {
	cp := c.Query("q")
	if cp == "" {
		c.JSON(http.StatusBadRequest, models.Error("q parameter required"))
		return
	}
	limit := parseLimit(c, 100)
	result, err := h.service.SearchByCodePostal(c.Request.Context(), cp, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}

func (h *Handler) SearchByDateCreation(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" {
		c.JSON(http.StatusBadRequest, models.Error("from parameter required"))
		return
	}
	limit := parseLimit(c, 100)
	result, err := h.service.SearchByDateCreation(c.Request.Context(), from, to, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}

func (h *Handler) SearchByCommune(c *gin.Context) {
	commune := c.Query("q")
	if commune == "" {
		c.JSON(http.StatusBadRequest, models.Error("q parameter required"))
		return
	}
	limit := parseLimit(c, 100)
	result, err := h.service.SearchByCommune(c.Request.Context(), commune, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}

func (h *Handler) SearchByEtatAdministratif(c *gin.Context) {
	etat := c.Query("q")
	if etat == "" {
		c.JSON(http.StatusBadRequest, models.Error("q parameter required"))
		return
	}
	limit := parseLimit(c, 100)
	result, err := h.service.SearchByEtatAdministratif(c.Request.Context(), etat, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}

func (h *Handler) SearchMultiCriteria(c *gin.Context) {
	criteria := models.CompanySearchCriteria{
		NafCode:           c.Query("naf"),
		Denomination:      c.Query("denomination"),
		CodePostal:        c.Query("codepostal"),
		Commune:           c.Query("commune"),
		EtatAdministratif: c.Query("etat"),
		DateCreationFrom:  c.Query("from"),
		DateCreationTo:    c.Query("to"),
	}
	limit := parseLimit(c, 100)
	result, err := h.service.SearchMultiCriteria(c.Request.Context(), criteria, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}
