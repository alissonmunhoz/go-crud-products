package service

import (
	"github.com/alissonmunhoz/go-crud-products/internal/config"
	"gorm.io/gorm"
)

var (
	logger *config.Logger
	db     *gorm.DB
)

func InitializeHandler() {
	logger = config.GetLogger("handler")
	db = config.GetMySQL()
}
