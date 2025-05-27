package handlers

import (
	"csv-importer/api"
	"fmt"
	"os"
)

func HandleAPI() {
	api.StartAPIServer()
	fmt.Println("🚀 API Server started")
	os.Exit(0)
}
