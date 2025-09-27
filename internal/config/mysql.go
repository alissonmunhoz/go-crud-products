package config

import (
	"fmt"
	"os"
	"time"

	"github.com/alissonmunhoz/go-crud-products/internal/schemas"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeMySQL() (*gorm.DB, error) {
	logger := GetLogger("mysql")

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	pass := getEnv("DB_PASSWORD", "root")
	name := getEnv("DB_NAME", "products")

	// DSN recomendado pelo GORM: charset utf8mb4 + parseTime + loc
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("mysql connection error: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	// Migrações
	if err := db.AutoMigrate(&schemas.Product{}); err != nil {
		logger.Errorf("mysql automigration error: %v", err)
		return nil, err
	}

	return db, nil
}

// helper para env com default
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
