package api

import (
	"csv-importer/api/handlers"
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

	// CORS
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

	server := &Server{db: db, router: router}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "CSV Importer API"})
	})

	// Tables
	api.GET("/tables", handlers.ListTables(s.db))
	api.GET("/tables/structure", handlers.GetCompleteStructure(s.db))
	api.GET("/tables/:name/info", handlers.GetTableInfo(s.db))
	api.GET("/tables/:name/columns", handlers.GetTableColumns(s.db))

	// Data
	api.GET("/data/:table/preview", handlers.PreviewTable(s.db))
	api.GET("/data/:table/values/:column", handlers.GetColumnValues(s.db))

	// Search
	api.GET("/search/:table/:column", handlers.SearchTable(s.db))
	api.GET("/count/:table/:column", handlers.CountRows(s.db))
	api.GET("/export/:table/:column", handlers.ExportData(s.db))
}

func (s *Server) Start(port string) error {
	log.Printf("🚀 API Server starting on port %s", port)
	log.Printf("📡 Health: http://localhost%s/api/health", port)
	log.Printf("📊 Tables: http://localhost%s/api/tables", port)
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
