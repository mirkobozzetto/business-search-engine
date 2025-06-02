package api

import (
	"csv-importer/api/middleware"
	"csv-importer/api/services/company"
	"csv-importer/api/services/data"
	"csv-importer/api/services/export"
	"csv-importer/api/services/search"
	"csv-importer/api/services/tables"
	"csv-importer/config"
	"csv-importer/database"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
)

type Server struct {
	db     *sql.DB
	router *gin.Engine
	logger *slog.Logger

	dataHandler    *data.Handler
	searchHandler  *search.Handler
	tableHandler   *tables.Handler
	exportHandler  *export.Handler
	companyHandler *company.Handler
}

func createLogger() *slog.Logger {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	})
	return slog.New(handler)
}

func NewServer(db *sql.DB) *Server {
	logger := createLogger()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.Use(middleware.ResponseMiddleware())

	dataService := data.NewDataService(db)
	dataHandler := data.NewHandler(dataService)

	searchService := search.NewSearchService(db)
	searchHandler := search.NewHandler(searchService)

	tableService := tables.NewTableService(db)
	tableHandler := tables.NewHandler(tableService)

	exportService := export.NewExportService(db)
	exportHandler := export.NewHandler(exportService)

	companyService := company.NewCompanyService(db)
	companyHandler := company.NewHandler(companyService)

	server := &Server{
		db:             db,
		router:         router,
		logger:         logger,
		dataHandler:    dataHandler,
		searchHandler:  searchHandler,
		tableHandler:   tableHandler,
		exportHandler:  exportHandler,
		companyHandler: companyHandler,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)
		responseHelper.Success(gin.H{"status": "ok", "message": "CSV Importer API"})
	})

	tablesGroup := api.Group("/tables")
	{
		tablesGroup.GET("", s.tableHandler.ListTables())
		tablesGroup.GET("/structure", s.tableHandler.GetCompleteStructure())

		tablesGroup.Use(middleware.ValidateTableName())
		tablesGroup.GET("/:name/info", s.tableHandler.GetTableInfo())
		tablesGroup.GET("/:name/columns", s.tableHandler.GetTableColumns())
	}

	dataGroup := api.Group("/data")
	dataGroup.Use(middleware.ValidateTableName())
	{
		dataGroup.GET("/:table/preview",
			middleware.ParseLimitParam(5, 100),
			s.dataHandler.PreviewTable(),
		)

		dataGroup.Use(middleware.ValidateColumnName(s.db))
		dataGroup.GET("/:table/values/:column",
			middleware.ParseLimitParam(20, 1000),
			s.dataHandler.GetColumnValues(),
		)
	}

	searchGroup := api.Group("/search")
	searchGroup.Use(middleware.ValidateTableName())
	{
		searchGroup.GET("/:table/:column",
			middleware.ValidateSearchQuery(),
			middleware.ParseLimitParam(50, 1000),
			s.searchHandler.SearchTable(),
		)

		searchGroup.GET("/:table/multi",
			middleware.ValidateSearchQuery(),
			middleware.ParseLimitParam(50, 1000),
			s.searchHandler.SearchMultipleColumns(),
		)
	}

	api.GET("/search/nacecode", s.searchHandler.SearchNaceCode())

	countGroup := api.Group("/count")
	countGroup.Use(middleware.ValidateTableName())
	countGroup.Use(middleware.ValidateColumnName(s.db))
	{
		countGroup.GET("/:table/:column",
			middleware.ValidateSearchQuery(),
			s.searchHandler.CountRows(),
		)
	}

	exportGroup := api.Group("/export")
	exportGroup.Use(middleware.ValidateTableName())
	{
		exportGroup.GET("/:table",
			middleware.ParseLimitParam(10000, 100000),
			middleware.ParseFormatParam(),
			s.exportHandler.ExportData(),
		)
	}

	companyGroup := api.Group("/companies")
	{
		companyGroup.GET("/search/nace", s.companyHandler.SearchByNaceCode())
		companyGroup.GET("/search/denomination", s.companyHandler.SearchByDenomination())
		companyGroup.GET("/search/zipcode", s.companyHandler.SearchByZipcode())
	}
}

func (s *Server) Start(port string) error {
	s.logger.Info("🚀 API Server starting",
		slog.String("port", port),
	)
	s.logger.Info("📡 Health endpoint",
		slog.String("url", "http://localhost"+port+"/api/health"),
	)
	s.logger.Info("📊 Tables structure",
		slog.String("url", "http://localhost"+port+"/api/tables/structure"),
	)
	s.logger.Info("🔍 NACE search",
		slog.String("url", "http://localhost"+port+"/api/search/nacecode"),
	)
	s.logger.Info("🔍 Company search",
		slog.String("url", "http://localhost"+port+"/api/companies/search/nace"),
	)
	s.logger.Info("🔍 Company search",
		slog.String("url", "http://localhost"+port+"/api/companies/search/denomination"),
	)

	return s.router.Run(port)
}

func StartAPIServer() {
	logger := createLogger()
	slog.SetDefault(logger)

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("❌ Database connection failed",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database connection",
				slog.String("error", err.Error()),
			)
		}
	}()

	server := NewServer(db)
	if err := server.Start(":8080"); err != nil {
		slog.Error("❌ Server failed to start",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
