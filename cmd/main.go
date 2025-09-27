package main

import (
	"github.com/alissonmunhoz/go-crud-products/internal/config"
	"github.com/alissonmunhoz/go-crud-products/internal/router"
)

var (
	logger config.Logger
)

// @title           Products API
// @version         1.0

// @host      localhost:8080
// @BasePath  /v1
// @schemes   http
func main() {

	logger = *config.GetLogger("main")

	err := config.Init()
	if err != nil {
		logger.Errorf("Config initalization error: %v", err)
		return
	}

	router.Initialize()
}
