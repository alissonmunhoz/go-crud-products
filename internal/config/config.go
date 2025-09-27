package config

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	logger *Logger
)

func Init() error {
	var err error

	db, err = InitializeMySQL()

	if err != nil {
		return fmt.Errorf("error initializing mysql: %v", err)
	}

	return nil
}

func GetMySQL() *gorm.DB {
	return db
}

func GetLogger(p string) *Logger {

	logger = NewLogger(p)
	return logger
}
