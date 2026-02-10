package api

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"sirene-importer/config"
	"sirene-importer/database"
	"sirene-importer/api/services/company"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
)

type Server struct {
	db             *sql.DB
	router         *gin.Engine
	logger         *slog.Logger
	companyHandler *company.Handler
}

func StartAPIServer() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("DB connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	server := NewServer(db)
	server.Run(":8081")
}

func NewServer(db *sql.DB) *Server {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	}))
	slog.SetDefault(logger)
	companyService := company.NewCompanyService(db)
	companyHandler := company.NewHandler(companyService)
	s := &Server{
		db:             db,
		router:         gin.Default(),
		logger:         logger,
		companyHandler: companyHandler,
	}
	s.setupRoutes()
	return s
}

func (s *Server) Run(addr string) {
	slog.Info("SIRENE France API", "addr", addr)
	s.router.Run(addr)
}

func (s *Server) setupRoutes() {
	s.router.Use(corsMiddleware())
	api := s.router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "sirene-france"})
	})
	companies := api.Group("/companies")
	companies.GET("/search/naf", s.companyHandler.SearchByNafCode)
	companies.GET("/search/denomination", s.companyHandler.SearchByDenomination)
	companies.GET("/search/codepostal", s.companyHandler.SearchByCodePostal)
	companies.GET("/search/commune", s.companyHandler.SearchByCommune)
	companies.GET("/search/etatadministratif", s.companyHandler.SearchByEtatAdministratif)
	companies.GET("/search/datecreation", s.companyHandler.SearchByDateCreation)
	companies.GET("/search/multi", s.companyHandler.SearchMultiCriteria)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
