package api

import (
	"csv-importer/api/handlers"
	"csv-importer/api/middleware"
	"csv-importer/config"
	"csv-importer/database"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	db     *sql.DB
	router *gin.Engine
}

func NewServer(db *sql.DB) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Global middlewares
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

	server := &Server{db: db, router: router}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		responseHelper := middleware.GetResponseHelper(c)
		responseHelper.Success(gin.H{"status": "ok", "message": "CSV Importer API"})
	})

	// Tables routes with middleware validation
	tablesGroup := api.Group("/tables")
	{
		tablesGroup.GET("", handlers.ListTables(s.db))
		tablesGroup.GET("/structure", handlers.GetCompleteStructure(s.db))

		// Routes requiring table validation
		tablesGroup.Use(middleware.ValidateTableName())
		tablesGroup.GET("/:name/info", handlers.GetTableInfo(s.db))
		tablesGroup.GET("/:name/columns", handlers.GetTableColumns(s.db))
	}

	// Data routes with middleware
	dataGroup := api.Group("/data")
	dataGroup.Use(middleware.ValidateTableName())
	{
		dataGroup.GET("/:table/preview",
			middleware.ParseLimitParam(5, 100),
			handlers.PreviewTable(s.db),
		)

		dataGroup.Use(middleware.ValidateColumnName(s.db))
		dataGroup.GET("/:table/values/:column",
			middleware.ParseLimitParam(20, 1000),
			handlers.GetColumnValues(s.db),
		)
	}

	// Search routes with validation middleware
	searchGroup := api.Group("/search")
	searchGroup.Use(middleware.ValidateTableName())
	searchGroup.Use(middleware.ValidateColumnName(s.db))
	{
		searchGroup.GET("/:table/:column",
			middleware.ValidateSearchQuery(),
			middleware.ParseLimitParam(50, 1000),
			handlers.SearchTable(s.db),
		)
	}

	// Count routes
	countGroup := api.Group("/count")
	countGroup.Use(middleware.ValidateTableName())
	countGroup.Use(middleware.ValidateColumnName(s.db))
	{
		countGroup.GET("/:table/:column",
			middleware.ValidateSearchQuery(),
			handlers.CountRows(s.db),
		)
	}

	// Export routes with enhanced middleware
	exportGroup := api.Group("/export")
	exportGroup.Use(middleware.ValidateTableName())
	{
		exportGroup.GET("/:table",
			middleware.ParseLimitParam(10000, 100000),
			middleware.ParseFormatParam(),
			handlers.ExportData(s.db),
		)
	}
}

func (s *Server) Start(port string) error {
	log.Printf("🚀 API Server starting on port %s", port)
	log.Printf("📡 Health: http://localhost%s/api/health", port)
	log.Printf("📊 Tables: http://localhost%s/api/tables", port)
	log.Printf("📊 Tables Structure: http://localhost%s/api/tables/structure", port)
	return s.router.Run(port)
}

func StartAPIServer() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}
	defer db.Close()

	server := NewServer(db)
	if err := server.Start(":8080"); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}
