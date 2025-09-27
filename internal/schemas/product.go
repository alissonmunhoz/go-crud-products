package schemas

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string
	Price       int64
	Quantity    int32
	Description string
}

type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Price       int64     `json:"price"`
	Quantity    int32     `json:"quantity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt,omitempty"`
}
