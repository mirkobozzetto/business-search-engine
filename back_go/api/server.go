package api

import (
	datahandlers "csv-importer/api/data/handlers"
	exporthandlers "csv-importer/api/export/handlers"
	"csv-importer/api/middleware"
	searchhandlers "csv-importer/api/search/handlers"
	tableshandlers "csv-importer/api/tables/handlers"
	"csv-importer/config"
	"csv-importer/database"
	"database/sql"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

type Server struct {
	db     *sql.DB
	router *gin.Engine
	logger *slog.Logger
}

func NewServer(db *sql.DB) *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	}))

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

	server := &Server{
		db:     db,
		router: router,
		logger: logger,
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
		tablesGroup.GET("", tableshandlers.ListTables(s.db))
		tablesGroup.GET("/structure", tableshandlers.GetCompleteStructure(s.db))

		tablesGroup.Use(middleware.ValidateTableName())
		tablesGroup.GET("/:name/info", tableshandlers.GetTableInfo(s.db))
		tablesGroup.GET("/:name/columns", tableshandlers.GetTableColumns(s.db))
	}

	dataGroup := api.Group("/data")
	dataGroup.Use(middleware.ValidateTableName())
	{
		dataGroup.GET("/:table/preview",
			middleware.ParseLimitParam(5, 100),
			datahandlers.PreviewTable(s.db),
		)

		dataGroup.Use(middleware.ValidateColumnName(s.db))
		dataGroup.GET("/:table/values/:column",
			middleware.ParseLimitParam(20, 1000),
			datahandlers.GetColumnValues(s.db),
		)
	}

	searchGroup := api.Group("/search")
	searchGroup.Use(middleware.ValidateTableName())
	searchGroup.Use(middleware.ValidateColumnName(s.db))
	{
		searchGroup.GET("/:table/:column",
			middleware.ValidateSearchQuery(),
			middleware.ParseLimitParam(50, 1000),
			searchhandlers.SearchTable(s.db),
		)
	}

	api.GET("/search/nacecode", searchhandlers.SearchNaceCode(s.db))

	countGroup := api.Group("/count")
	countGroup.Use(middleware.ValidateTableName())
	countGroup.Use(middleware.ValidateColumnName(s.db))
	{
		countGroup.GET("/:table/:column",
			middleware.ValidateSearchQuery(),
			searchhandlers.CountRows(s.db),
		)
	}

	exportGroup := api.Group("/export")
	exportGroup.Use(middleware.ValidateTableName())
	{
		exportGroup.GET("/:table",
			middleware.ParseLimitParam(10000, 100000),
			middleware.ParseFormatParam(),
			exporthandlers.ExportData(s.db),
		)
	}
}



func (s *Server) Start(port string) error {
	s.logger.Info("API Server starting",
		slog.String("port", port),
		slog.String("health_endpoint", "http://localhost"+port+"/api/health"),
		slog.String("tables_structure", "http://localhost"+port+"/api/tables/structure"),
		slog.String("nace_search", "http://localhost"+port+"/api/search/nacecode"),
	)

	return s.router.Run(port)
}

func StartAPIServer() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("Database connection failed",
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
		slog.Error("Server failed to start",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
