package company

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sirene-importer/api/models"
	"strconv"
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

func parseOffset(c *gin.Context) int {
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			return parsed
		}
	}
	return 0
}

func (h *Handler) SearchByNafCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, models.Error("code parameter required"))
		return
	}
	limit := parseLimit(c, 100)
	offset := parseOffset(c)
	result, err := h.service.SearchByNafCode(c.Request.Context(), code, limit, offset)
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
	offset := parseOffset(c)
	result, err := h.service.SearchByDenomination(c.Request.Context(), query, limit, offset)
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
	offset := parseOffset(c)
	result, err := h.service.SearchByCodePostal(c.Request.Context(), cp, limit, offset)
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
	offset := parseOffset(c)
	result, err := h.service.SearchByDateCreation(c.Request.Context(), from, to, limit, offset)
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
	offset := parseOffset(c)
	result, err := h.service.SearchByCommune(c.Request.Context(), commune, limit, offset)
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
	offset := parseOffset(c)
	result, err := h.service.SearchByEtatAdministratif(c.Request.Context(), etat, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}

func (h *Handler) SearchMultiCriteria(c *gin.Context) {
	criteria := models.CompanySearchCriteria{
		NafCode:            c.Query("naf"),
		Denomination:       c.Query("denomination"),
		CodePostal:         c.Query("codepostal"),
		Commune:            c.Query("commune"),
		EtatAdministratif:  c.Query("etat"),
		DateCreationFrom:   c.Query("from"),
		DateCreationTo:     c.Query("to"),
		CategorieJuridique: c.Query("categorie_juridique"),
		TrancheEffectifs:   c.Query("tranche_effectifs"),
	}
	limit := parseLimit(c, 100)
	offset := parseOffset(c)
	result, err := h.service.SearchMultiCriteria(c.Request.Context(), criteria, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}

func (h *Handler) SearchByIdentifier(c *gin.Context) {
	identifier := c.Param("identifier")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, models.Error("identifier parameter required (SIREN 9 digits or SIRET 14 digits)"))
		return
	}
	result, err := h.service.SearchByIdentifier(c.Request.Context(), identifier)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.Success(result))
}
