// handler/product_mapper.go
package service

import (
	"time"

	"github.com/alissonmunhoz/go-crud-products/internal/schemas"
	"gorm.io/gorm"
)

func toProductResponse(p schemas.Product) schemas.ProductResponse {
	var del *time.Time

	if da, ok := any(p.DeletedAt).(gorm.DeletedAt); ok {
		if da.Valid {
			t := da.Time
			del = &t
		}
	} else if tptr, ok := any(p.DeletedAt).(*time.Time); ok {
		del = tptr
	}

	return schemas.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Quantity:    p.Quantity,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt: func() time.Time {
			if del != nil {
				return *del
			}
			return time.Time{}
		}(),
	}
}
