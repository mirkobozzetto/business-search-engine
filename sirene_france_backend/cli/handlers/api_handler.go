package handlers

import (
	"fmt"
	"sirene-importer/api"
)

func HandleAPI() {
	fmt.Println("Starting SIRENE France API server...")
	api.StartAPIServer()
}
